package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		dir := "testdata/env"
		env, err := ReadDir(dir)
		require.NoError(t, err)

		require.Equal(t, "bar", env["BAR"].Value)
		require.Equal(t, "", env["EMPTY"].Value)
		require.Equal(t, "   foo\nwith new line", env["FOO"].Value)
		require.Equal(t, "\"hello\"", env["HELLO"].Value)
		require.Equal(t, "", env["UNSET"].Value)

		require.False(t, env["BAR"].NeedRemove)
		require.False(t, env["EMPTY"].NeedRemove)
		require.False(t, env["FOO"].NeedRemove)
		require.False(t, env["HELLO"].NeedRemove)
		require.True(t, env["UNSET"].NeedRemove)

		require.Equal(t, 5, len(env))
	})
	t.Run("dir not exists", func(t *testing.T) {
		dir := "testdata/1"
		_, err := ReadDir(dir)
		require.Error(t, err)
	})
	t.Run("not include dir", func(t *testing.T) {
		dir := "testdata"
		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 1, len(env))
	})
}
