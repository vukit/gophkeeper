package service_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vukit/gophkeeper/internal/client/model"
	"github.com/vukit/gophkeeper/internal/client/service"
)

func TestCryptoArray(t *testing.T) {
	tests := []struct {
		name    string
		user    model.User
		message string
	}{
		{
			name:    "case 1",
			user:    model.User{Username: "Mark", Password: "superSecret"},
			message: "This is test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, err := service.NewCryptoService(&tt.user)
			require.Nil(t, err)
			encrypted := cs.Encrypt([]byte(tt.message))
			decrypted, err := cs.Decrypt(encrypted)
			require.Nil(t, err)
			assert.Equal(t, tt.message, string(decrypted))
		})
	}

}

func TestCryptoFile(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "gophkeeper_crypto_test")
	os.RemoveAll(dir)

	err := os.Mkdir(dir, os.ModePerm)
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	src, err := ioutil.TempFile(dir, "*")
	assert.Nil(t, err)
	_, err = io.CopyN(src, rand.Reader, 1024)
	defer src.Close()

	cs, err := service.NewCryptoService(&model.User{Username: "mark", Password: "superSecret"})
	assert.Nil(t, err)

	encyptedDst, err := cs.EncryptFile(src)
	assert.Nil(t, err)

	dst, err := cs.DecryptFile(io.NopCloser(encyptedDst))
	assert.Nil(t, err)

	srcBytes, err := ioutil.ReadAll(src)
	assert.Nil(t, err)

	dstBytes, err := ioutil.ReadAll(dst)
	assert.Nil(t, err)

	assert.True(t, bytes.Equal(srcBytes, dstBytes))
}
