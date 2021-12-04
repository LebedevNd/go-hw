package main

import (
	"os"
)

func main() {
	directory := os.Args[1]
	args := []string{os.Args[2]}
	args = append(args, os.Args[3:]...)

	env, err := ReadDir(directory)
	if err != nil {
		return
	}

	RunCmd(args, env)
}
