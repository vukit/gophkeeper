package model

import (
	"errors"
	"fmt"
	"strings"
)

// User модель пользователя приложения
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const maxUserUsernameLegth = 64

var (
	ErrUserUsernamePasswordEmpity = errors.New("username and/or password empity")
	ErrUserLongUsername           = fmt.Errorf("username length is more than %d characters", maxUserUsernameLegth)
)

// Validate проверяет корректность модели пользователя приложения
func (r *User) Validate() error {
	if strings.TrimSpace(r.Username) == "" || strings.TrimSpace(r.Password) == "" {
		return ErrUserUsernamePasswordEmpity
	}

	if len(r.Username) > maxUserUsernameLegth {
		return ErrUserLongUsername
	}

	return nil
}
