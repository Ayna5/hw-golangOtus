package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		fmt.Println("cmd length must be greater than 0")
		return 1
	}
	for key, val := range env {
		err := os.Setenv(key, val.Value)
		if err != nil {
			return 1
		}
		fmt.Println("key: ", key, "value: ", val)
	}

	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Env = os.Environ()
	if err := c.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return 1
	}
	return 0
}
