package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const (
	RunCmdErrorWrongParametersCount = iota + 10
	RunCmdErrorPrepareEnvironment
	RunCmdErrorStartCommand
	RunCmdErrorCommand
)

type RunCmdError interface {
	error
	Code() int
}

type RunError struct {
	err  error
	code int
}

func (err *RunError) Error() string {
	return err.err.Error()
}

func (err *RunError) Code() int {
	return err.code
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, ioErr io.Writer, ioOut io.Writer, reader io.Reader) (int, RunCmdError) { //nolint: lll
	if len(cmd) == 0 {
		return 0, &RunError{
			err:  fmt.Errorf("wrong command parameters count"),
			code: RunCmdErrorWrongParametersCount,
		}
	}
	var envErr error
	command := exec.Command(cmd[0], cmd[1:]...) //nolint: gosec
	command.Stderr = ioErr
	command.Stdin = reader
	command.Stdout = ioOut
	command.Env, envErr = prepareEnvs(env)
	if envErr != nil {
		return 0, &RunError{
			err:  fmt.Errorf("prepare environment variables: %w", envErr),
			code: RunCmdErrorPrepareEnvironment,
		}
	}
	startErr := command.Start()
	if startErr != nil {
		return 0, &RunError{
			err:  fmt.Errorf("command start: %w", startErr),
			code: RunCmdErrorStartCommand,
		}
	}

	err := command.Wait()

	var exitError *exec.ExitError

	if err != nil && !errors.As(err, &exitError) {
		return command.ProcessState.ExitCode(), &RunError{
			err:  fmt.Errorf("command wait: %w", err),
			code: RunCmdErrorCommand,
		}
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
