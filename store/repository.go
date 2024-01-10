package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ito0804takuya/go_todo_app/clock"
	"github.com/ito0804takuya/go_todo_app/config"
	"github.com/jmoiron/sqlx"
)

// レポジトリの生成（MySQLへの接続）
func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	// MySQLへの接続
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		),
	)
	if err != nil {
		return nil, nil, err
	}

	// タイムアウト設定
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// 接続テスト
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}

	// database/sqlをラップしたsqlxを使う（レコードから構造体へのマッピングが楽なので）
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { _ = db.Close() }, nil
}

type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

// 更新系
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// 参照系
type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

// 各インターフェースが型に沿っているか確認する（ビルドエラーで分かるように）
var (
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer = (*sqlx.DB)(nil)
	_ Execer = (*sqlx.DB)(nil)
	_ Execer = (*sqlx.Tx)(nil)
)

type Repository struct {
	Clocker clock.Clocker
}

