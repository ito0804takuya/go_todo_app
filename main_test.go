package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	// contextを用意
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// テスト対象のrun関数を実行
		return run(ctx)
	})

	// リクエスト
	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// 期待する出力結果
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// run関数に終了通知を送信
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
