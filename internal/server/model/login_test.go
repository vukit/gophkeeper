package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukit/gophkeeper/internal/server/model"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		want     error
	}{
		{
			name:     "case 1",
			username: "mark",
			password: "secret",
			want:     nil,
		},
		{
			name:     "case 2",
			username: "mark",
			password: "",
			want:     model.ErrLoginPasswordEmpity,
		},
		{
			name:     "case 3",
			username: "",
			password: "secret",
			want:     model.ErrLoginUsernameEmpity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			login := model.Login{Username: tt.username, Password: tt.password}
			assert.Equal(t, tt.want, login.Validate())
		})
	}

}
