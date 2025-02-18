package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
	"to-do-checklist/internal/database"
	"to-do-checklist/internal/kafka"
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

func GetHandler(c echo.Context, db *gorm.DB, producer *kafka.Producer) error {
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
	msg := fmt.Sprintf("Request received at %s, method: GET, query: %s", time.Now().Format(time.RFC3339), c.Request().RequestURI)
	if err := producer.Produce(msg, "tasks"); err != nil {
		return fmt.Errorf("There was error with Kafka Producer: %v", err)
	}
	return respondWithSuccess(c, http.StatusOK, tasks)
}

func PostTaskHandler(c echo.Context, db *gorm.DB, producer *kafka.Producer) error {
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

	msg := fmt.Sprintf("Request received at %s, method: POST, query: %s", time.Now().Format(time.RFC3339), c.Request().RequestURI)
	if err := producer.Produce(msg, "tasks"); err != nil {
		return fmt.Errorf("There was error with Kafka Producer: %v", err)
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("OK", "Task was added successfully"))
}

func PatchTaskHandler(c echo.Context, db *gorm.DB, producer *kafka.Producer) error {
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

	msg := fmt.Sprintf("Request received at %s, method: PATCH, query: %s", time.Now().Format(time.RFC3339), c.Request().RequestURI)
	if err := producer.Produce(msg, "tasks"); err != nil {
		return fmt.Errorf("There was error with Kafka Producer: %v", err)
	}

	return respondWithSuccess(c, http.StatusOK, createResponse("OK", "Task was updated successfully"))
}

func DeleteTaskHandler(c echo.Context, db *gorm.DB, producer *kafka.Producer) error {
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

	msg := fmt.Sprintf("Request received at %s, method: DELETE, query: %s", time.Now().Format(time.RFC3339), c.Request().RequestURI)
	if err := producer.Produce(msg, "tasks"); err != nil {
		return fmt.Errorf("There was error with Kafka Producer: %v", err)
	}
	
	return respondWithSuccess(c, http.StatusOK, createResponse("OK", "The message was deleted successfully"))
}
