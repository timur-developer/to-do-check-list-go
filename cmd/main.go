package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"to-do-checklist/internal/database"
	"to-do-checklist/internal/kafka"
	"to-do-checklist/internal/routes"
)

func main() {
	producer, err := kafka.NewProducer([]string{"kafka:9092"})
	if err != nil {
		fmt.Errorf("There was error with creating producer")
	}
	db, _ := database.InitDB()

	e := echo.New()
	routes.RegisterRoutes(e, db, producer)

	e.Logger.Fatal(e.Start(":8080"))

}
