package model

import (
	"errors"
	"strings"
)

// File модель файла пользователя
type File struct {
	ID       int    `json:"id"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	MetaInfo string `json:"metainfo"`
}

var ErrFileMetainfoEmpity = errors.New("metainfo empity")

// Validate проверяет корректность модели файла пользователя
func (r *File) Validate() error {
	if strings.TrimSpace(r.MetaInfo) == "" {
		return ErrFileMetainfoEmpity
	}

	return nil
}
