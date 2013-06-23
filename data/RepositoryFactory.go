package data

type RepositoryFactory interface {
	CreateUserRepository() UserRepository
	CreateCallbackRepository() CallbackRepository
}
