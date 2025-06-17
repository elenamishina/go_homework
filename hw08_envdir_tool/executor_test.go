package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		cmd := []string{}
		env := make(Environment)
		retCode := RunCmd(cmd, env)
		require.Equal(t, 1, retCode)
	})
	t.Run("run error", func(t *testing.T) {
		cmd := []string{"testdata/ech.sh"}
		env := make(Environment)
		retCode := RunCmd(cmd, env)
		require.Equal(t, 1, retCode)
	})
	t.Run("command error", func(t *testing.T) {
		cmd := []string{"ls", "path"}
		env := make(Environment)
		retCode := RunCmd(cmd, env)
		require.Equal(t, 2, retCode)
	})
	t.Run("simple without args", func(t *testing.T) {
		cmd := []string{"testdata/echo.sh"}
		env := make(Environment)
		env["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		retCode := RunCmd(cmd, env)
		require.Equal(t, 0, retCode)
	})
	t.Run("simple with args", func(t *testing.T) {
		cmd := []string{"echo", "hello"}
		env := make(Environment)
		retCode := RunCmd(cmd, env)
		require.Equal(t, 0, retCode)
	})
}
