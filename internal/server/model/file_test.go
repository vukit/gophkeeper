package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukit/gophkeeper/internal/server/model"
)

func TestFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		filename string
		metainfo string
		want     error
	}{
		{
			name:     "case 1",
			path:     "/tmp/test.pdf",
			filename: "test.pdf",
			metainfo: "this is test file",
			want:     nil,
		},
		{
			name:     "case 2",
			path:     "",
			filename: "test.pdf",
			metainfo: "this is test file",
			want:     model.ErrFilePathEmpity,
		},
		{
			name:     "case 3",
			path:     "/tmp/test.pdf",
			filename: "",
			metainfo: "this is test file",
			want:     model.ErrFileNameEmpity,
		},
		{
			name:     "case 4",
			path:     "/tmp/test.pdf",
			filename: "test.pdf",
			metainfo: "",
			want:     model.ErrFileMetainfoEmpity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := model.File{Path: tt.path, Name: tt.filename, MetaInfo: tt.metainfo}
			assert.Equal(t, tt.want, file.Validate())
		})
	}

}
