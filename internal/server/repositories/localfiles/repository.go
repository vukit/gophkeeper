package localfiles

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
)

// RepoLocalFiles структура репозитория локальных файлов
type RepoLocalFiles struct {
	path string
}

// NewRepo возвращает репозиторий локальных файлов
func NewRepo(path string) (repo *RepoLocalFiles, err error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	repo = &RepoLocalFiles{path: path}

	return repo, nil
}

// SaveFile сохраняет файл пользователя
func (repo RepoLocalFiles) SaveFile(ctx context.Context, src io.ReadCloser) (filePath string, err error) {
	filePath, err = repo.getPath()
	if err != nil {
		return "", err
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return filePath, err
}

// DeleteFile удаляет файл прользователя
func (repo RepoLocalFiles) DeleteFile(ctx context.Context, filePath string) (err error) {
	if filePath == "" {
		return nil
	}

	err = os.Remove(filePath)
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	return nil
}

// GetFile возвращает файл пользователя
func (repo RepoLocalFiles) GetFile(ctx context.Context, filePath string) (fileReader io.ReadCloser, err error) {
	return os.Open(filePath)
}

// Close закрывает файловый репозиторий
func (repo RepoLocalFiles) Close() error {
	return nil
}

func (repo RepoLocalFiles) getPath() (string, error) {
	for {
		bytes, err := generateRandomBytes(32)
		if err != nil {
			return "", err
		}

		targetDir := filepath.Join(repo.path, string(bytes[0:2]), string(bytes[2:4]), string(bytes[4:6]))
		targetFile := filepath.Join(targetDir, string(bytes[6:]))

		if _, err := os.Stat(targetFile); errors.Is(err, os.ErrNotExist) {
			if _, err := os.Stat(targetDir); errors.Is(err, os.ErrNotExist) {
				err := os.MkdirAll(targetDir, os.ModePerm)
				if err != nil {
					return "", nil
				}
			}

			return targetFile, nil
		}
	}
}

func generateRandomBytes(n int) ([]byte, error) {
	const letters = "0123456789abcdef"

	result := make([]byte, n)

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return nil, err
		}

		result[i] = letters[num.Int64()]
	}

	return result, nil
}
