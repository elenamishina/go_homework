package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	command.Env = generateEnv(env)

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
		return 1
	}
	return 0
}

func generateEnv(customEnv Environment) []string {
	var env []string

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if _, exists := customEnv[pair[0]]; !exists {
			env = append(env, e)
		}
	}

	for name, val := range customEnv {
		if !val.NeedRemove {
			env = append(env, name+"="+val.Value)
		}
	}
	return env
}
