package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Print("directory is not exist")
		os.Exit(1)
	}

	d := os.Args[1]
	env, err := ReadDir(d)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "can't get environment variables from dir"))
		os.Exit(1)
	}

	exitCode := RunCmd(os.Args[2:], env)
	os.Exit(exitCode)
}
