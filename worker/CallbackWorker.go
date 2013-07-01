package worker

import (
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
	"net/http"
	"sync"
	"time"
)

type CallbackWorker struct {
	callbackRepository data.CallbackRepository

	startStopLock sync.Mutex
	runningLock   sync.Mutex

	// Set to true by Start(), set to false at doWork() ending
	isRunning bool

	stop    chan bool
	stopped sync.WaitGroup
}

func NewCallbackWorker(callbackRepository data.CallbackRepository) *CallbackWorker {
	return &CallbackWorker{
		callbackRepository: callbackRepository,
		stop:               make(chan bool),
	}
}

func (worker *CallbackWorker) Start() {
	worker.startStopLock.Lock()
	defer worker.startStopLock.Unlock()

	if worker.isRunning {
		return
	}

	worker.isRunning = true
	go worker.doWork()
}

func (worker *CallbackWorker) Stop() {
	worker.startStopLock.Lock()
	defer worker.startStopLock.Unlock()

	if worker.isRunning {
		worker.stop <- true
	}
}

func (worker *CallbackWorker) doWork() {
	worker.runningLock.Lock()
	defer worker.stopped.Done()
	defer worker.runningLock.Unlock()
	var workers sync.WaitGroup

	for {
		select {
		case <-worker.stop:
			Log.Notice("Stopping worker")
		case <-time.After(50 * time.Millisecond):
			callback, err := worker.callbackRepository.GetNextAndBumpNextAttemptTimeStamp(1 * time.Minute)
			if err != nil {
				Log.Error("Error while getting next callback: %v", err.Error())
				continue
			} else {
				if callback == nil {
					// No callback available that needs to be called.
					continue
				}

				Log.Info("Handling callback %v\n\turl:%v", callback.Id, callback.Request.Url)
				workers.Add(1)
				go worker.executeCallback(callback, workers)
			}
		}
	}

	Log.Notice("Waiting for workers to finish")
	workers.Wait()
	Log.Info("Done!")

	worker.isRunning = false
}

func (worker *CallbackWorker) executeCallback(callback *model.Callback, workers sync.WaitGroup) {
	defer workers.Done()

	start := time.Now()
	response, err := http.Get(callback.Request.Url)
	stop := time.Now()
	responseTime := stop.Sub(start)
	if err != nil {
		Log.Error("Error while executing callback:\n\rerror:%s\n\rurl:%s\n\tresponse time:%s",
			err.Error(), callback.Request.Url, responseTime)
	}
	defer response.Body.Close()

	attempt := &model.CallbackAttempt{
		Id:        model.NewObjectId(),
		Timestamp: time.Now(),
		Success:   response.StatusCode == http.StatusOK,
		Response: &model.HttpResponseInfo{
			HttpStatusCode: response.StatusCode,
			HttpStatusText: response.Status,
			ResponseTime:   responseTime,
		},
	}

	Log.Info("New attempt added")
	if err := worker.callbackRepository.AddAttemptToCallback(callback.Id, attempt); err != nil {
		Log.Error("Error while saving callback attempt: ", err.Error())
	}
}
