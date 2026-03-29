package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

// addTaskHandler godoc
// @Summary     Добавить задачу
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       task body db.Task true "Параметры задачи"
// @Success     200 {object} map[string]string "id созданной задачи"
// @Failure     400 {object} map[string]string "ошибка валидации"
// @Router      /api/task [post]
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "ошибка десериализации JSON: "+err.Error())
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "ошибка сохранения задачи: "+err.Error())
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("дата указана в неверном формате, ожидается YYYYMMDD")
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}
