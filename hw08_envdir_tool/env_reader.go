package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, fileInfo := range files {
		name := fileInfo.Name()
		fileBuff := make([]byte, fileInfo.Size())

		file, _ := os.Open(dir + "/" + name)
		_, err = file.Read(fileBuff)
		if err != nil {
			return nil, err
		}

		fileLines := strings.Split(string(fileBuff), "\n")

		firstLine := strings.TrimRight(fileLines[0], " ")
		firstLine = string(bytes.ReplaceAll([]byte(firstLine), []byte{0x00}, []byte("\n")))
		env[name] = EnvValue{
			firstLine,
			len(firstLine) == 0,
		}
	}

	return env, err
}
