package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./go-envdir </path/to/env/dir> <command> [args...]")
		os.Exit(1)
	}

	envPath := os.Args[1]
	command := os.Args[2]
	args := os.Args[3:]

	env, err := ReadDir(envPath)
	if err != nil {
		fmt.Println("Error read env")
		os.Exit(1)
	}

	cmd := append([]string{command}, args...)
	retCode := RunCmd(cmd, env)
	os.Exit(retCode)
}
