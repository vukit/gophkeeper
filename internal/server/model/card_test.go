package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukit/gophkeeper/internal/server/model"
)

func TestCard(t *testing.T) {
	tests := []struct {
		name   string
		bank   string
		number string
		date   string
		cvv    string
		want   error
	}{
		{
			name:   "case 1",
			bank:   "Good bank",
			number: "4111111111111111",
			date:   "12/24",
			cvv:    "123",
			want:   nil,
		},
		{
			name:   "case 2",
			bank:   "",
			number: "4111111111111111",
			date:   "12/24",
			cvv:    "123",
			want:   model.ErrCardBankEmpity,
		},
		{
			name:   "case 3",
			bank:   "Good Bank",
			number: "",
			date:   "12/24",
			cvv:    "123",
			want:   model.ErrCardNumberEmpity,
		},
		{
			name:   "case 4",
			bank:   "Good bank",
			number: "41111111111111111",
			date:   "",
			cvv:    "123",
			want:   model.ErrCardDateEmpity,
		},
		{
			name:   "case 5",
			bank:   "Good bank",
			number: "4111111111111111",
			date:   "12/24",
			cvv:    "",
			want:   model.ErrCardCVVEmpity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := model.Card{Bank: tt.bank, Number: tt.number, Date: tt.date, CVV: tt.cvv}
			assert.Equal(t, tt.want, card.Validate())
		})
	}

}
