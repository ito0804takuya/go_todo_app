package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ito0804takuya/go_todo_app/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	// テストしやすいようにrun関数に切り出した
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// SIGINT（割り込みシグナル）かSIGTERM（終了シグナル）を受け取るとグレースフルシャットダウンするよう設定
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

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