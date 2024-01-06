package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	t.Skip("リファクタリング中")

	l, err := net.Listen("tcp", "localhost:0") // "0"指定 : 利用可能なポートを動的に選択してくれる
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}

	// contextを用意
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// テスト対象のrun関数を実行
		return run(ctx)
	})

	// リクエスト
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
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
