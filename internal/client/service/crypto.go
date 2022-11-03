package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"

	"github.com/vukit/gophkeeper/internal/client/model"
)

// CryptoService структура сервиса симметричного шифрования
type CryptoService struct {
	aesgcm cipher.AEAD
	nonce  []byte
}

// NewCryptoService возвращает сервис симметричного шифрования для пользователя приложения
func NewCryptoService(user *model.User) (cs *CryptoService, err error) {
	cs = &CryptoService{}

	key := sha256.Sum256([]byte(user.Password))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	cs.aesgcm, err = cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	cs.nonce = key[len(key)-cs.aesgcm.NonceSize():]

	return cs, nil
}

// Encrypt шифрует масссив src
func (r *CryptoService) Encrypt(src []byte) []byte {
	return r.aesgcm.Seal(nil, r.nonce, src, nil)
}

// Decrypt расшифровывает масссив src
func (r *CryptoService) Decrypt(src []byte) ([]byte, error) {
	dst, err := r.aesgcm.Open(nil, r.nonce, src, nil)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// EncryptFile шифрует файл src
func (r *CryptoService) EncryptFile(src io.ReadCloser) (dst io.Reader, err error) {
	inBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	outBytes := r.Encrypt(inBytes)

	dst = bytes.NewReader(outBytes)

	return dst, nil
}

// DecryptFile шифрует файл src
func (r *CryptoService) DecryptFile(src io.ReadCloser) (dst io.Reader, err error) {
	inBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	outBytes, err := r.Decrypt(inBytes)
	if err != nil {
		return nil, err
	}

	dst = bytes.NewReader(outBytes)

	return dst, nil
}
