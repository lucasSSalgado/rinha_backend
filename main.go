package main

import (
	"rinha/controller"
	"rinha/infra"
	"rinha/repository"
	"rinha/service"

	"github.com/patrickmn/go-cache"
)

func main() {
	db := infra.CreateConnection()
	c := cache.New(cache.NoExpiration, cache.NoExpiration)

	clientRepo := repository.NewClientRepository(db)
	clientServ := service.NewClientService(clientRepo, c)
	cont := controller.NewControl(clientServ)

	cont.InitRoutes()
}
