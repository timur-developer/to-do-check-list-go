package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"to-do-checklist/internal/database"
	"to-do-checklist/internal/kafka"
	"to-do-checklist/internal/routes"
)

func main() {
	address := []string{"kafka:9092"}
	producer, err := kafka.NewProducer(address)
	if err != nil {
		log.Errorf("There was error with creating producer: %v", err)
	}
	consumer, err := kafka.NewConsumer(address, "tasks", "my-group")
	if err != nil {
		log.Errorf("There was error with creating consumer: %v", err)
	}

	go consumer.Start()
	defer consumer.Stop()

	db, _ := database.InitDB()

	e := echo.New()
	routes.RegisterRoutes(e, db, producer)

	e.Logger.Fatal(e.Start(":8080"))
}
