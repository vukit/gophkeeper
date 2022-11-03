package localfiles_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vukit/gophkeeper/internal/server/repositories/localfiles"
)

func TestSaveLocalFile(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "gophkeeper_local_files")
	defer os.RemoveAll(dir)

	ctx := context.Background()
	repoFile, err := localfiles.NewRepo(dir)
	defer repoFile.Close()

	srcFile, err := ioutil.TempFile(dir, "src_")
	require.Nil(t, err)

	dstFilePath, err := repoFile.SaveFile(ctx, srcFile)
	require.Nil(t, err)

	_, err = os.Stat(dstFilePath)

	require.Nil(t, err)
}

func TestDeleteLocalFile(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "gophkeeper_local_files")
	defer os.RemoveAll(dir)

	ctx := context.Background()
	repoFile, err := localfiles.NewRepo(dir)
	defer repoFile.Close()

	srcFile, err := ioutil.TempFile(dir, "src_")
	require.Nil(t, err)

	dstFilePath, err := repoFile.SaveFile(ctx, srcFile)
	require.Nil(t, err)

	_, err = os.Stat(dstFilePath)
	require.Nil(t, err)

	err = repoFile.DeleteFile(ctx, dstFilePath)
	require.Nil(t, err)

	_, err = os.Stat(dstFilePath)
	assert.ErrorIs(t, err, os.ErrNotExist)

}

func TestGetLocalFile(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "gophkeeper_local_files")
	defer os.RemoveAll(dir)

	ctx := context.Background()
	repoFile, err := localfiles.NewRepo(dir)
	defer repoFile.Close()

	srcFile, err := ioutil.TempFile(dir, "src_")
	require.Nil(t, err)

	dstFilePath, err := repoFile.SaveFile(ctx, srcFile)
	require.Nil(t, err)

	_, err = os.Stat(dstFilePath)
	require.Nil(t, err)

	fileReader, err := repoFile.GetFile(ctx, dstFilePath)
	assert.Nil(t, err)
	fileReader.Close()
}
