package main

import (
	"github.com/labstack/echo/v4"
	"to-do-checklist/internal/database"
	"to-do-checklist/internal/routes"
)

func main() {

	db, _ := database.InitDB()

	e := echo.New()
	routes.RegisterRoutes(e, db)

	e.Logger.Fatal(e.Start(":8080"))

}
