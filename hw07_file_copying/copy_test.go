package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("oversize offset case", func(t *testing.T) {
		err := Copy("testdata/test_input.txt", "testdata/dgad.txt", 1000000, 0)
		require.Error(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("oversize offset case", func(t *testing.T) {
		err := Copy("testdata/unsup_file", "testdata/dgad.txt", 1000000, 0)
		require.Error(t, ErrUnsupportedFile, err)
	})

	t.Run("creeate file case", func(t *testing.T) {
		_ = Copy("testdata/test_input.txt", "testdata/dgad.txt", 0, 0)
		require.FileExistsf(t, "testdata/dgad.txt", "hmm")
	})

	_ = os.Remove("testdata/dgad.txt")
}
