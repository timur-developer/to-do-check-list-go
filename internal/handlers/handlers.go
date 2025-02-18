package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"to-do-checklist/internal/database"
)

func respondWithError(c echo.Context, status int, text string) error {
	return c.JSON(status, map[string]string{"error": text})
}

func respondWithSuccess(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, data)
}

func createResponse(status, task string) map[string]string {
	return map[string]string{
		"status": status,
		"task":   task,
	}
}

func GetHandler(c echo.Context, db *gorm.DB) error {
	var tasks []database.Task
	query := db.Model(database.Task{})
	searchText := c.QueryParam("searchText")
	params := c.QueryParams()

	if importance, exists := params["importance"]; exists {
		query = query.Where("importance = ?", importance[0])
	}
	if doneParam, exists := params["is_done"]; exists {
		query = query.Where("is_done = ?", doneParam[0])
	}
	if searchText != "" {
		query = query.Where("(task_name ILIKE ? OR task_description ILIKE ?)", "%"+searchText+"%", "%"+searchText+"%")
	}

	if err := query.Find(&tasks).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Error while getting tasks")
	}

	if len(tasks) == 0 {
		return respondWithError(c, http.StatusNotFound, "There are no any tasks")
	}
	
	return respondWithSuccess(c, http.StatusOK, tasks)
}

func PostTaskHandler(c echo.Context, db *gorm.DB) error {
	var task database.Task
	if err := c.Bind(&task); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect POST request")
	}

	// Добавление задачи с ID, большим на 1 ID, максимального из всех задач в таблице
	//var maxID int
	//db.Model(&database.Task{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
	//task.ID = maxID + 1

	if err := db.Create(&task).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not add the task")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("OK", "Task was added successfully"))
}

func PatchTaskHandler(c echo.Context, db *gorm.DB) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect ID")
	}

	var task database.Task
	if err := db.First(&task, id).Error; err != nil {
		return respondWithError(c, http.StatusNotFound, fmt.Sprintf("There is no task with id: %d", id))
	}

	var updatedTask database.Task
	if err := c.Bind(&updatedTask); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect PATCH request")
	}

	if err := db.Model(&task).Updates(updatedTask).Error; err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Could not update the task")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("OK", "Task was updated successfully"))
}

func DeleteTaskHandler(c echo.Context, db *gorm.DB) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Incorrect ID")
	}

	var task database.Task
	if err := db.First(&task, id).Error; err != nil {
		return respondWithError(c, http.StatusNotFound, "Task was not found")
	}

	if err := db.Delete(&task).Error; err != nil {
		return respondWithError(c, http.StatusBadRequest, "Could not delete the message")
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("OK", "The message was deleted successfully"))
}
