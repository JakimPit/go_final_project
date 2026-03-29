package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"go_final_project/pkg/db"
)

// getTaskHandler godoc
// @Summary     Получить задачу по ID
// @Tags        tasks
// @Produce     json
// @Param       id query string true "ID задачи"
// @Success     200 {object} db.Task
// @Failure     400 {object} map[string]string "ошибка"
// @Router      /api/task [get]
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, err.Error(), http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, task, http.StatusOK)
}

// updateTaskHandler godoc
// @Summary     Обновить задачу
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       task body db.Task true "Обновлённые параметры задачи (включая id)"
// @Success     200 {object} map[string]string "пустой объект при успехе"
// @Failure     400 {object} map[string]string "ошибка валидации"
// @Router      /api/task [put]
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, err.Error(), http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}

// deleteTaskHandler godoc
// @Summary     Удалить задачу
// @Tags        tasks
// @Produce     json
// @Param       id query string true "ID задачи"
// @Success     200 {object} map[string]string "пустой объект при успехе"
// @Failure     400 {object} map[string]string "ошибка"
// @Router      /api/task [delete]
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, err.Error(), http.StatusNotFound)
		} else {
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}
