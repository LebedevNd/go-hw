package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	setEnv(env)

	commandString := cmd[0]
	command := exec.Command(commandString, cmd[1:]...)
	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		println(err.Error())
		return 0
	}

	return 1
}

func setEnv(env Environment) {
	for name, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(name)
			if err != nil {
				return
			}
		} else {
			err := os.Setenv(name, value.Value)
			if err != nil {
				return
			}
		}
	}
}
