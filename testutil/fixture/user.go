package fixture

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/ito0804takuya/go_todo_app/entity"
)

// fixture = テストで必要となるダミーデータ

func User(u *entity.User) *entity.User {
	result := &entity.User{
		ID:       entity.UserID(rand.Int()),
		Name:     "hoge" + strconv.Itoa(rand.Int())[:5],
		Password: "passWord",
		Role:     "admin",
		Created:  time.Now(),
		Modified: time.Now(),
	}

	if u == nil {
		return result
	}

	// uが持っているフィールドがあれば、それに置き換えてresultを返す
	if u.ID != 0 {
		result.ID = u.ID
	}
	if u.Name != "" {
		result.Name = u.Name
	}
	if u.Password != "" {
		result.Password = u.Password
	}
	if u.Role != "" {
		result.Role = u.Role
	}
	if u.Created.IsZero() {
		result.Created = u.Created
	}
	if u.Modified.IsZero() {
		result.Modified = u.Modified
	}

	return result
}
