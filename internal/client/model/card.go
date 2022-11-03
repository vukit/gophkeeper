package model

import (
	"errors"
	"fmt"
	"strings"
)

// Card модель банковской карты приложения
type Card struct {
	ID       int    `json:"id"`
	Bank     string `json:"bank"`
	Number   string `json:"number"`
	Date     string `json:"date"`
	CVV      string `json:"cvv"`
	MetaInfo string `json:"metainfo"`
}

const (
	maxCardNumberLegth = 19
	maxCardCVVLegth    = 3
)

var (
	ErrCardBankEmpity     = errors.New("bank empity")
	ErrCardNumberEmpity   = errors.New("number empity")
	ErrCardLongNumber     = fmt.Errorf("number length is more than %d characters", maxCardNumberLegth)
	ErrCardDateEmpity     = errors.New("date empity")
	ErrCardCVVEmpity      = errors.New("cvv empity")
	ErrCardLongCVV        = fmt.Errorf("cvv length is more than %d characters", maxCardCVVLegth)
	ErrCardMetainfoEmpity = errors.New("metainfo empity")
)

// Validate проверяет корректность модели банковской карты приложения
func (r *Card) Validate() error {
	if strings.TrimSpace(r.Bank) == "" {
		return ErrCardBankEmpity
	}

	if strings.TrimSpace(r.Number) == "" {
		return ErrCardNumberEmpity
	}

	if len(r.Number) > maxCardNumberLegth {
		return ErrCardLongNumber
	}

	if strings.TrimSpace(r.Date) == "" {
		return ErrCardDateEmpity
	}

	if strings.TrimSpace(r.CVV) == "" {
		return ErrCardCVVEmpity
	}

	if len(r.CVV) > maxCardCVVLegth {
		return ErrCardLongCVV
	}

	if strings.TrimSpace(r.MetaInfo) == "" {
		return ErrCardMetainfoEmpity
	}

	return nil
}
