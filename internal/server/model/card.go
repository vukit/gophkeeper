package model

import (
	"errors"
	"strings"
)

// Card модель банковской карты сервера
type Card struct {
	ID       int    `json:"id"`
	UserID   int    `json:"-"`
	Bank     string `json:"bank"`
	Number   string `json:"number"`
	Date     string `json:"date"`
	CVV      string `json:"cvv"`
	MetaInfo string `json:"metainfo"`
}

var (
	ErrCardBankEmpity   = errors.New("bank empity")
	ErrCardNumberEmpity = errors.New("number empity")
	ErrCardDateEmpity   = errors.New("date empity")
	ErrCardCVVEmpity    = errors.New("cvv empity")
)

// Validate проверяет корректность модели банковской карты сервера
func (r *Card) Validate() error {
	if strings.TrimSpace(r.Bank) == "" {
		return ErrCardBankEmpity
	}

	if strings.TrimSpace(r.Number) == "" {
		return ErrCardNumberEmpity
	}

	if strings.TrimSpace(r.Date) == "" {
		return ErrCardDateEmpity
	}

	if strings.TrimSpace(r.CVV) == "" {
		return ErrCardCVVEmpity
	}

	return nil
}
