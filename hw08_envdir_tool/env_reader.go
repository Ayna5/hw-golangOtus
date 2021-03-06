package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
//nolint:cyclop
func ReadDir(dir string) (Environment, error) {
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read directory for path: %s", dir)
	}

	res := make(map[string]EnvValue)
	for _, file := range fileInfo {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), "=") {
			continue
		}
		openFile, err := os.OpenFile(filepath.Join(dir, file.Name()), os.O_RDONLY, os.ModeDir)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot open file for path: %s", file.Name())
		}

		r := bufio.NewReader(openFile)
		l, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				res[file.Name()] = EnvValue{Value: "", NeedRemove: true}
				openFile.Close()
				continue
			}
			openFile.Close()
			return nil, fmt.Errorf("can't read line from file %v: %w", file.Name(), err)
		}

		str := bytes.TrimRight(l, " \t")
		str1 := bytes.ReplaceAll(str, []byte("\x00"), []byte("\n"))
		if len(str1) == 0 {
			delete(res, file.Name())
			res[file.Name()] = EnvValue{NeedRemove: true}
			openFile.Close()
			continue
		}
		res[file.Name()] = EnvValue{Value: string(str1)}
	}
	return res, nil
}
