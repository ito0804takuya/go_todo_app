package auth

import (
	"bytes"
	"context"
	"testing"

	"github.com/ito0804takuya/go_todo_app/entity"
	"github.com/ito0804takuya/go_todo_app/testutil/fixture"
)

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("want %s, but got %s", want, rawPubKey)
	}
	want = []byte("-----BEGIN RSA PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("want %s, but got %s", want, rawPrivKey)
	}
}

func TestJWTer_GenerateToken(t *testing.T) {
	moq := &StoreMock{}

	wantID := entity.UserID(20)
	u := fixture.User(&entity.User{ID: wantID})

	// GenerateTokenメソッドの中で使うSaveメソッドをmoq版に置き換え
	moq.SaveFunc = func(ctx context.Context, key string, userID  entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, but got %d", wantID, userID)
		}
		return nil
	}

	sut, err := NewJWTer(moq)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	got, err := sut.GenerateToken(ctx, *u)
	if err != nil {
		t.Fatalf("not want err:%v", err)
	}
	if len(got) == 0 {
		t.Errorf("token is empty")
	}
}
