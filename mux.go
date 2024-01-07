package main

import "net/http"

// muxはマルチプレクサのことで、複数の入力を単一の出力に結合するための装置のこと。

func NewMux() http.Handler {
	mux := http.NewServeMux()

	// ヘルスチェック用
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json: charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	return mux
}