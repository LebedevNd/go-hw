package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("success run", func(t *testing.T) {
		returnCode := RunCmd([]string{"/bin/bash"}, Environment{})
		require.Equal(t, 1, returnCode)
	})

	t.Run("wrong cmd run", func(t *testing.T) {
		returnCode := RunCmd([]string{"/something/wrong"}, Environment{})
		require.Equal(t, 0, returnCode)
	})
}
