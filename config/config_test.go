package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	// 起動ポートを設定
	wantPort := 3333
	t.Setenv("PORT", fmt.Sprint(wantPort))
	// 環境変数から読み込んだconfigオブジェクトを生成
	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}

	// PORTについてテスト
	if got.Port != wantPort {
		t.Errorf("want %d, but %d", wantPort, got.Port)
	}

	// ENVについてテスト
	wantEnv := "dev"
	if got.Env != wantEnv {
		t.Errorf("want %s, but %s", wantEnv, got.Env)
	}
}