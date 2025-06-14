package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	testInput := "testdata/input.txt"
	testOutput := "out.txt"
	t.Run("source is system file", func(t *testing.T) {
		fileFrom := "/dev/urandom"
		fileTo := testOutput
		var offset int64
		var limit int64

		err := Copy(fileFrom, fileTo, offset, limit)
		require.Equal(t, ErrUnsupportedFile, err)
	})
	t.Run("offset less than zero", func(t *testing.T) {
		fileFrom := testInput
		fileTo := testOutput
		var offset int64 = -1
		var limit int64

		err := Copy(fileFrom, fileTo, offset, limit)
		require.Equal(t, ErrOffsetNegative, err)
	})
	t.Run("limit less than zero", func(t *testing.T) {
		fileFrom := testInput
		fileTo := testOutput
		var offset int64
		var limit int64 = -1

		err := Copy(fileFrom, fileTo, offset, limit)
		require.Equal(t, ErrLimitNegative, err)
	})
	t.Run("source file not exist", func(t *testing.T) {
		fileFrom := "input.txt"
		fileTo := testOutput
		var offset int64
		var limit int64

		err := Copy(fileFrom, fileTo, offset, limit)
		require.NotEqual(t, nil, err)
	})
	t.Run("offset exceeds file size", func(t *testing.T) {
		fileFrom := testInput
		fileTo := testOutput
		var offset int64 = 10000
		var limit int64

		err := Copy(fileFrom, fileTo, offset, limit)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("full file", func(t *testing.T) {
		fileFrom := testInput
		tmpFile, _ := os.CreateTemp("", testOutput)
		tmpFile.Close()
		var offset int64
		var limit int64

		err := Copy(fileFrom, tmpFile.Name(), offset, limit)
		require.Equal(t, nil, err)
		os.Remove(tmpFile.Name())
	})

	t.Run("set limit and ofset", func(t *testing.T) {
		fileFrom := testInput
		tmpFile, _ := os.CreateTemp("", testOutput)
		tmpFile.Close()
		var offset int64 = 100
		var limit int64 = 1000

		err := Copy(fileFrom, tmpFile.Name(), offset, limit)
		require.Equal(t, nil, err)
		os.Remove(tmpFile.Name())
	})
	t.Run("set limit with offset exceeds file size", func(t *testing.T) {
		fileFrom := testInput
		tmpFile, _ := os.CreateTemp("", testOutput)
		tmpFile.Close()
		var offset int64 = 6000
		var limit int64 = 1000

		err := Copy(fileFrom, tmpFile.Name(), offset, limit)
		require.Equal(t, nil, err)
		os.Remove(tmpFile.Name())
	})
	t.Run("limit exceeds file size", func(t *testing.T) {
		fileFrom := testInput
		tmpFile, _ := os.CreateTemp("", testOutput)
		tmpFile.Close()
		var offset int64
		var limit int64 = 10000

		err := Copy(fileFrom, tmpFile.Name(), offset, limit)
		require.Equal(t, nil, err)
		os.Remove(tmpFile.Name())
	})
	t.Run("Destination file not exist", func(t *testing.T) {
		fileFrom := "testdata/"
		fileTo := testOutput
		var offset int64
		var limit int64

		err := Copy(fileFrom, fileTo, offset, limit)
		require.NotEqual(t, nil, err)
	})
	t.Run("Destination file not exist", func(t *testing.T) {
		fileFrom := "testdata/1/input.txt"
		fileTo := testOutput
		var offset int64
		var limit int64

		err := Copy(fileFrom, fileTo, offset, limit)
		require.NotEqual(t, nil, err)
	})
}
