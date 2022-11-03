package model

import (
	"errors"
	"strings"
)

// File модель файла пользователя
type File struct {
	ID       int    `json:"id"`
	UserID   int    `json:"-"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	MetaInfo string `json:"metainfo"`
}

var (
	ErrFilePathEmpity     = errors.New("file path empity")
	ErrFileNameEmpity     = errors.New("name empity")
	ErrFileMetainfoEmpity = errors.New("metainfo empity")
)

// Validate проверяет корректность модели файла пользователя
func (r *File) Validate() error {
	if strings.TrimSpace(r.Path) == "" {
		return ErrFilePathEmpity
	}

	if strings.TrimSpace(r.Name) == "" {
		return ErrFileNameEmpity
	}

	if strings.TrimSpace(r.MetaInfo) == "" {
		return ErrFileMetainfoEmpity
	}

	return nil
}
