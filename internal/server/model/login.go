package model

import (
	"errors"
	"strings"
)

// Login модель логина сервера
type Login struct {
	ID       int    `json:"id"`
	UserID   int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	MetaInfo string `json:"metainfo"`
}

var (
	ErrLoginUsernameEmpity = errors.New("username empity")
	ErrLoginPasswordEmpity = errors.New("password empity")
)

// Validate проверяет корректность модели логина сервера
func (r *Login) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return ErrLoginUsernameEmpity
	}

	if strings.TrimSpace(r.Password) == "" {
		return ErrLoginPasswordEmpity
	}

	return nil
}
