package main

import (
	"rinha/controller"
	"rinha/infra"
	"rinha/repository"
	"rinha/service"
)

func main() {
	db := infra.CreateConnection()

	clientRepo := repository.NewClientRepository(db)
	clientServ := service.NewClientService(clientRepo)
	cont := controller.NewControl(clientServ)

	cont.InitRoutes()
}
