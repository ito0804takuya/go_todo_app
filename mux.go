package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/ito0804takuya/go_todo_app/handler"
	"github.com/ito0804takuya/go_todo_app/store"
)

// muxはマルチプレクサのことで、複数の入力を単一の出力に結合するための装置のこと。

func NewMux() http.Handler {
	mux := chi.NewRouter()

	// ヘルスチェック用
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json: charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	// AddTask
	v := validator.New()
	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks", at.ServerHTTP)

	// ListTask
	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", lt.ServeHTTP)

	return mux
}