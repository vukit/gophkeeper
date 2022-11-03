package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukit/gophkeeper/internal/client/model"
)

func TestUser(t *testing.T) {
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
			want:     model.ErrUserUsernamePasswordEmpity,
		},
		{
			name:     "case 3",
			username: "",
			password: "secret",
			want:     model.ErrUserUsernamePasswordEmpity,
		},
		{
			name:     "case 4",
			username: "",
			password: "",
			want:     model.ErrUserUsernamePasswordEmpity,
		},
		{
			name:     "case 5",
			username: "longusernamelongusernamelongusernamelongusernamelongusernamelongusername",
			password: "secret",
			want:     model.ErrUserLongUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := model.User{Username: tt.username, Password: tt.password}
			assert.Equal(t, tt.want, client.Validate())
		})
	}

}
