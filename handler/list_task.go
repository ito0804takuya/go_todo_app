package handler

import (
	"net/http"

	"github.com/ito0804takuya/go_todo_app/entity"
	"github.com/ito0804takuya/go_todo_app/store"
	"github.com/jmoiron/sqlx"
)

type ListTask struct {
	// Store *store.TaskStore
	DB *sqlx.DB
	Repo store.Repository
}

type task struct {
	ID entity.TaskID `json:"id"`
	Title string `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// tasks := lt.Store.All()
	tasks, err := lt.Repo.ListTasks(ctx, lt.DB)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}
	rsp := []task{}
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID: t.ID, 
			Title: t.Title, 
			Status: t.Status,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}