package gocqrs

// import (
//   "github.com/pjvds/go-cqrs/sourcing"
//   "github.com/pjvds/go-cqrs/storage"
// )

// type CallbackRepository struct {
//   repository storage.Repository
// }

// func NewCallbackRepository(backend storage.RepositoryBackend) *CallbackRepository{
//   repository := storage.NewRepository(backend)

//   return &CallbackRepository{
//     repository: repository
//   }
// }

// func (r *CallbackRepository) Add(callback *model.Callback) error {
//   source := sourcing.GetState(callback)
//   return r.repository.Add(source)
// }

// func (r *CallbackRepository) Get(model.ObjectId) (*model.Callback, error) {
//   callback := new(Callback)
//   callback.EventSource = sourcing.AttachNew(callback)

//   r.repository.Get(sourceId, callback)
// }
