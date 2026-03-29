package api

import (
	"encoding/json"
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
		writeError(w, "не указан идентификатор задачи")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, task)
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
		writeError(w, "ошибка десериализации JSON: "+err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, "не указан идентификатор задачи")
		return
	}

	if task.Title == "" {
		writeError(w, "не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, struct{}{})
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
		writeError(w, "не указан идентификатор задачи")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, struct{}{})
}
