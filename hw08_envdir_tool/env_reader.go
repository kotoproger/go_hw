package main

import (
	"bufio"
	"fmt"
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
	result := make(Environment)

	dirData, readError := os.ReadDir(dir)
	if readError != nil {
		return nil, fmt.Errorf("read dir: %w", readError)
	}
	for _, dirEntry := range dirData {
		if dirEntry.IsDir() {
			continue
		}
		name := dirEntry.Name()
		if !isValidEnvName(name) {
			continue
		}

		value, readFileError := readFirstLine(dir + "/" + name)
		value = strings.ReplaceAll(value, "\x00", "\n")
		if readFileError != nil {
			return nil, fmt.Errorf("read dir entry: %w", readFileError)
		}
		if value == "" {
			result[name] = EnvValue{NeedRemove: true}
		} else {
			result[name] = EnvValue{Value: value}
		}
	}
	return result, nil
}

func readFirstLine(name string) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		err = scanner.Err()
		if err != nil {
			return "", fmt.Errorf("scan file: %w", err)
		}

		return "", nil
	}

	return scanner.Text(), nil
}

func isValidEnvName(name string) bool {
	return !strings.Contains(name, "=")
}
