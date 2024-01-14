package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/ito0804takuya/go_todo_app/clock"
	"github.com/ito0804takuya/go_todo_app/config"
	"github.com/ito0804takuya/go_todo_app/handler"
	"github.com/ito0804takuya/go_todo_app/store"
)

// muxはマルチプレクサのことで、複数の入力を単一の出力に結合するための装置のこと。

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()

	// ヘルスチェック用
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json: charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	// AddTask
	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := store.Repository{ Clocker: clock.RealClocker{} }
	at := &handler.AddTask{DB: db, Repo: &r, Validator: v}
	mux.Post("/tasks", at.ServerHTTP)

	// ListTask
	lt := &handler.ListTask{DB: db, Repo: r}
	mux.Get("/tasks", lt.ServeHTTP)

	return mux, cleanup, nil
}