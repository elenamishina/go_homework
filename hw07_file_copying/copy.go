package main

import (
	"errors"
	"io"
	"os"
	"strings"

	//nolint:depguard
	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetNegative        = errors.New("offset less than zero")
	ErrLimitNegative         = errors.New("limit less than zero")
	ErrUnsupportedStat       = errors.New("unsupported stat")
)

func isSystemFile(file string) bool {
	systemDirs := []string{"/dev", "/proc", "/sys", "/run"}
	for _, systemDir := range systemDirs {
		if strings.HasPrefix(file, systemDir) {
			return true
		}
	}
	return false
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if isSystemFile(fromPath) {
		return ErrUnsupportedFile
	}
	if offset < 0 {
		return ErrOffsetNegative
	}
	if limit < 0 {
		return ErrLimitNegative
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	infoSrcFile, err := srcFile.Stat()
	if err != nil {
		return ErrUnsupportedStat
	}
	if offset > infoSrcFile.Size() {
		return ErrOffsetExceedsFileSize
	}

	ableSizeCopy := infoSrcFile.Size() - offset
	if limit == 0 || limit > ableSizeCopy {
		limit = ableSizeCopy
	}

	if offset > 0 {
		_, err = srcFile.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	progressBar := pb.Full.Start64(limit)
	progressBar.Set(pb.Bytes, true)
	progressBar.Set(pb.SIBytesPrefix, true)
	proxyReader := progressBar.NewProxyReader(srcFile)

	_, err = io.CopyN(dstFile, proxyReader, limit)
	progressBar.Finish()
	if err != nil {
		return err
	}
	return nil
}
