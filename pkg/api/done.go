package api

import (
	"errors"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

// doneHandler godoc
// @Summary     Отметить задачу выполненной
// @Tags        tasks
// @Produce     json
// @Param       id query string true "ID задачи"
// @Success     200 {object} map[string]string "пустой объект при успехе"
// @Failure     400 {object} map[string]string "ошибка"
// @Router      /api/task/done [post]
func doneHandler(w http.ResponseWriter, r *http.Request) {
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

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.UpdateDate(id, next); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}
