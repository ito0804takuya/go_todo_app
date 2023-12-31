package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	// ポート番号
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]

	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}

	// テストしやすいようにrun関数に切り出した
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	// 別ゴルーチンでHTTPサーバーを起動
	eg.Go(func() error {
		// ErrServerClosedで意図的に終了した場合を除く
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close:rub %+v", err)
			return err
		}
		return nil
	})

	// キャンセル通知を待機
	<-ctx.Done()
	// キャンセル通知を受け取ったら終了
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// eg.Go()で起動した、errgroupのすべてのゴルーチンが完了するのを待ち、その結果(error)を返す
	return eg.Wait()
}