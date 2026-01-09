package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, ioErr io.Writer, ioOut io.Writer, reader io.Reader) (returnCode int, err error) { //nolint: lll
	if len(cmd) == 0 {
		return 0, fmt.Errorf("wrong command parameters count")
	}
	var envErr error
	command := exec.Command(cmd[0], cmd[1:]...) //nolint: gosec
	command.Stderr = ioErr
	command.Stdin = reader
	command.Stdout = ioOut
	command.Env, envErr = prepareEnvs(env)
	if envErr != nil {
		return 0, fmt.Errorf("prepare environment variables: %w", envErr)
	}
	err = command.Start()
	if err != nil {
		return 0, fmt.Errorf("command start: %w", err)
	}

	err = command.Wait()

	var exitError *exec.ExitError

	if err != nil && !errors.As(err, &exitError) {
		return command.ProcessState.ExitCode(), fmt.Errorf("command wait: %w", err)
	}

	return command.ProcessState.ExitCode(), nil
}

func prepareEnvs(env Environment) ([]string, error) {
	for key, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return nil, fmt.Errorf("unset environment variable: %w", err)
			}

			continue
		}

		err := os.Setenv(key, value.Value)
		if err != nil {
			return nil, fmt.Errorf("unset environment variable: %w", err)
		}
	}

	return os.Environ(), nil
}
