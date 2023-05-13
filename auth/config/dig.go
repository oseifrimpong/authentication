package config

import (
	"auth-service/auth/repository"
	"auth-service/auth/service"

	"go.uber.org/dig"
)

var container = dig.New()

func BuildProject() *dig.Container {

	container.Provide(InitializeDB)
	container.Provide(Logger)

	container.Provide(repository.NewUserRepo)
	container.Provide(repository.NewVerificationRepo)

	container.Provide(service.NewUserService)
	container.Provide(service.NewVerificationService)
	container.Provide(service.NewKeyService)

	return container
}

func Invoke(i interface{}) error {
	return container.Invoke(i)
}
