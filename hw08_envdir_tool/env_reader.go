package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func getFirstLine(data []byte) []byte {
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return data[:i]
	}
	return data
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)
	entriesDir, err := os.ReadDir(dir)
	if err != nil {
		return env, err
	}

	for _, e := range entriesDir {
		if e.IsDir() {
			continue
		}
		if strings.Contains(e.Name(), "=") {
			fmt.Printf("file %s not read, invalid name", e.Name())
			continue
		}

		contentFile, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return env, err
		}
		firstLine := getFirstLine(contentFile)
		firstLine = bytes.ReplaceAll(firstLine, []byte{0x00}, []byte{'\n'})
		envValue := strings.TrimRight(string(firstLine), " \t")

		needRemove := false
		if len(contentFile) == 0 {
			needRemove = true
		}

		env[e.Name()] = EnvValue{Value: envValue, NeedRemove: needRemove}
	}
	return env, nil
}
