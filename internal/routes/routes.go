package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"to-do-checklist/internal/handlers"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	e.GET("/tasks", func(c echo.Context) error {
		return handlers.GetHandler(c, db)
	})
	e.POST("/create", func(c echo.Context) error {
		return handlers.PostTaskHandler(c, db)
	})
	e.PATCH("/edit/:id", func(c echo.Context) error {
		return handlers.PatchTaskHandler(c, db)
	})
	e.DELETE("/delete/:id", func(c echo.Context) error {
		return handlers.DeleteTaskHandler(c, db)
	})
}
