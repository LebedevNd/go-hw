package main

import (
	"errors"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileFrom, _ := os.Open(fromPath)
	fromSize, err := fileSize(fileFrom)

	buf := make([]byte, fromSize)

	if err != nil {
		return ErrUnsupportedFile
	}

	if fromSize < offset {
		return ErrOffsetExceedsFileSize
	}

	_, err = fileFrom.Read(buf)
	if err != nil {
		return err
	}

	err = fileFrom.Close()
	if err != nil {
		return err
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}

	if limit <= 0 || (limit+offset) > fromSize {
		limit = fromSize - offset
	}

	bar := pb.StartNew(int(limit))
	for i := 0; i < int(limit); i++ {
		_, err = fileTo.Write(buf[offset+int64(i) : offset+int64(i)+1])
		if err != nil {
			return err
		}
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.Finish()

	if err != nil {
		return err
	}

	err = fileTo.Close()
	if err != nil {
		return err
	}

	return nil
}

func fileSize(fileFrom *os.File) (int64, error) {
	statFrom, err := fileFrom.Stat()
	if err != nil {
		return 0, ErrUnsupportedFile
	}
	return statFrom.Size(), nil
}
