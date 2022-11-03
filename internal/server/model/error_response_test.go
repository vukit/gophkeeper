package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukit/gophkeeper/internal/server/model"
)

func TestError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "case 1",
			message: "this is test message",
			want:    "{\"error\": \"this is test message\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			error := model.ErrorResponse{Error: tt.message}
			assert.Equal(t, tt.want, error.String())
		})
	}

}
