package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukit/gophkeeper/internal/client/model"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		metainfo string
		want     error
	}{
		{
			name:     "case 1",
			username: "mark",
			password: "secret",
			metainfo: "secret message",
			want:     nil,
		},
		{
			name:     "case 2",
			username: "mark",
			password: "",
			metainfo: "secret message",
			want:     model.ErrLoginPasswordEmpity,
		},
		{
			name:     "case 3",
			username: "",
			password: "secret",
			metainfo: "secret message",
			want:     model.ErrLoginUsernameEmpity,
		},
		{
			name:     "case 4",
			username: "mark",
			password: "secret",
			metainfo: "",
			want:     model.ErrLoginMetaInfoEmpity,
		},
		{
			name:     "case 5",
			username: "longusernamelongusernamelongusernamelongusernamelongusernamelongusername",
			password: "secret",
			metainfo: "secret message",
			want:     model.ErrLoginLongUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			login := model.Login{Username: tt.username, Password: tt.password, MetaInfo: tt.metainfo}
			assert.Equal(t, tt.want, login.Validate())
		})
	}

}
