package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler godoc
// @Summary     Список задач
// @Tags        tasks
// @Produce     json
// @Param       search query string false "Поиск по тексту или дате (02.01.2006)"
// @Success     200 {object} tasksResp
// @Failure     500 {object} map[string]string "ошибка сервера"
// @Router      /api/tasks [get]
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	var (
		tasks []*db.Task
		err   error
	)

	if search == "" {
		tasks, err = db.Tasks(50)
	} else if t, parseErr := time.Parse("02.01.2006", search); parseErr == nil {
		tasks, err = db.TasksByDate(t.Format(dateFormat), 50)
	} else {
		tasks, err = db.TasksBySearch(search, 50)
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, tasksResp{Tasks: tasks}, http.StatusOK)
}
