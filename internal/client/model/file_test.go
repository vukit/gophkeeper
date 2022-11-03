package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukit/gophkeeper/internal/client/model"
)

func TestFile(t *testing.T) {
	tests := []struct {
		name     string
		metainfo string
		want     error
	}{
		{
			name:     "case 1",
			metainfo: "this is test file",
			want:     nil,
		},
		{
			name:     "case 2",
			metainfo: "",
			want:     model.ErrFileMetainfoEmpity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := model.File{MetaInfo: tt.metainfo}
			assert.Equal(t, tt.want, file.Validate())
		})
	}

}
