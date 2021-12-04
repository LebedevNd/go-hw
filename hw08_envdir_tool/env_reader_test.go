package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("wrong path", func(t *testing.T) {
		_, err := ReadDir("testdata/wrong_env")
		require.Error(t, err)
	})

	t.Run("test deleting second line", func(t *testing.T) {
		env, _ := ReadDir("testdata/env")
		require.NotNil(t, env["BAR"])
		require.Equal(t, env["BAR"].Value, "bar")
		require.False(t, env["BAR"].NeedRemove)
	})
}
