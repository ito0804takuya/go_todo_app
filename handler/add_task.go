package handler

import (
	"encoding/json"
	"net/http"
	// "time"

	"github.com/go-playground/validator"
	"github.com/ito0804takuya/go_todo_app/entity"
	"github.com/ito0804takuya/go_todo_app/store"
	"github.com/jmoiron/sqlx"
)

type AddTask struct {
	// Store *store.TaskStore
	DB *sqlx.DB
	Repo *store.Repository
	Validator *validator.Validate
}

func (at *AddTask) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// バリデーション実行
	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		Title: b.Title,
		Status: entity.TaskStatusTodo,
		// Created: time.Now(),
	}
	// id, err := store.Tasks.Add(t)
	err = at.Repo.AddTask(ctx, at.DB, t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID int `json:"id"`
	}{ID: int(t.ID)}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
