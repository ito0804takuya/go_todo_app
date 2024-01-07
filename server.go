package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	// SIGINT（割り込みシグナル）かSIGTERM（終了シグナル）を受け取るとグレースフルシャットダウンするよう設定
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)

	// 別ゴルーチンでHTTPサーバーを起動
	eg.Go(func() error {
		// ErrServerClosedで意図的に終了した場合を除く
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close:rub %+v", err)
			return err
		}
		return nil
	})

	// キャンセル通知を待機
	<-ctx.Done()
	// キャンセル通知を受け取ったら終了
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// eg.Go()で起動した、errgroupのすべてのゴルーチンが完了するのを待ち、その結果(error)を返す
	return eg.Wait()
}
