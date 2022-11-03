package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukit/gophkeeper/internal/client/model"
)

func TestCard(t *testing.T) {
	tests := []struct {
		name     string
		bank     string
		number   string
		date     string
		cvv      string
		metainfo string
		want     error
	}{
		{
			name:     "case 1",
			bank:     "Good bank",
			number:   "4111111111111111",
			date:     "12/24",
			cvv:      "123",
			metainfo: "This is test message",
			want:     nil,
		},
		{
			name:     "case 2",
			bank:     "",
			number:   "4111111111111111",
			date:     "12/24",
			cvv:      "123",
			metainfo: "This is test message",
			want:     model.ErrCardBankEmpity,
		},
		{
			name:     "case 3",
			bank:     "Good Bank",
			number:   "",
			date:     "12/24",
			cvv:      "123",
			metainfo: "This is test message",
			want:     model.ErrCardNumberEmpity,
		},
		{
			name:     "case 4",
			bank:     "Good bank",
			number:   "41111111111111111111",
			date:     "12/24",
			cvv:      "123",
			metainfo: "This is test message",
			want:     model.ErrCardLongNumber,
		},
		{
			name:     "case 5",
			bank:     "Good bank",
			number:   "41111111111111111",
			date:     "",
			cvv:      "123",
			metainfo: "This is test message",
			want:     model.ErrCardDateEmpity,
		},
		{
			name:     "case 6",
			bank:     "Good bank",
			number:   "4111111111111111",
			date:     "12/24",
			cvv:      "",
			metainfo: "This is test message",
			want:     model.ErrCardCVVEmpity,
		},
		{
			name:     "case 7",
			bank:     "Good bank",
			number:   "4111111111111111",
			date:     "12/24",
			cvv:      "1234",
			metainfo: "This is test message",
			want:     model.ErrCardLongCVV,
		},
		{
			name:     "case 8",
			bank:     "Good bank",
			number:   "4111111111111111",
			date:     "12/24",
			cvv:      "123",
			metainfo: "",
			want:     model.ErrCardMetainfoEmpity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := model.Card{Bank: tt.bank, Number: tt.number, Date: tt.date, CVV: tt.cvv, MetaInfo: tt.metainfo}
			assert.Equal(t, tt.want, card.Validate())
		})
	}

}
