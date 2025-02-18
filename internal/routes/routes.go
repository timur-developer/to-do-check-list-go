package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"to-do-checklist/internal/handlers"
	"to-do-checklist/internal/kafka"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, producer *kafka.Producer) {
	e.GET("/tasks", func(c echo.Context) error {
		return handlers.GetHandler(c, db, producer)
	})
	e.POST("/create", func(c echo.Context) error {
		return handlers.PostTaskHandler(c, db, producer)
	})
	e.PATCH("/edit/:id", func(c echo.Context) error {
		return handlers.PatchTaskHandler(c, db, producer)
	})
	e.DELETE("/delete/:id", func(c echo.Context) error {
		return handlers.DeleteTaskHandler(c, db, producer)
	})
}
