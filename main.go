package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"club/handler"
	"club/subscriber"

	club "club/proto/club"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("DMS.SMS.v1.service.club"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	club.RegisterClubHandler(service.Server(), new(handler.Club))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("DMS.SMS.v1.service.club", service.Server(), new(subscriber.Club))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
