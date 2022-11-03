package model

import (
	"errors"
	"fmt"
	"strings"
)

// Login модель логина приложения
type Login struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	MetaInfo string `json:"metainfo"`
}

const maxLoginUsernameLegth = 64

var (
	ErrLoginUsernameEmpity = errors.New("username empity")
	ErrLoginPasswordEmpity = errors.New("password empity")
	ErrLoginMetaInfoEmpity = errors.New("metainfo empity")
	ErrLoginLongUsername   = fmt.Errorf("username length is more than %d characters", maxLoginUsernameLegth)
)

// Validate проверяет корректность модели логина приложения
func (r *Login) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return ErrLoginUsernameEmpity
	}

	if strings.TrimSpace(r.Password) == "" {
		return ErrLoginPasswordEmpity
	}

	if len(r.Username) > maxLoginUsernameLegth {
		return ErrLoginLongUsername
	}

	if strings.TrimSpace(r.MetaInfo) == "" {
		return ErrLoginMetaInfoEmpity
	}

	return nil
}
