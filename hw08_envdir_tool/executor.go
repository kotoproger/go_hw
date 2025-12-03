package main

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, ioErr io.Writer, ioOut io.Writer, reader io.Reader, osSignals <-chan syscall.Signal) (returnCode int, err error) {
	if len(cmd) == 0 {
		return 0, fmt.Errorf("Wrong command parameters count")
	}
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stderr = ioErr
	command.Stdin = reader
	command.Stdout = ioOut
	command.Env = composeEnvSlice(env)
	err = command.Start()
	if err != nil {
		return 0, fmt.Errorf("command start: %w", err)
	}

	// не работает передача сигнала (((
	//sigChan := make(chan struct{})
	//go func(closer <-chan struct{}, osSignal <-chan syscall.Signal) {
	//	select {
	//	case <-closer:
	//		return
	//	case signal := <-osSignal:
	//		err := syscall.Kill(command.Process.Pid, signal
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}
	//}(sigChan, osSignals)

	err = command.Wait()
	//close(sigChan)

	var exitError *exec.ExitError

	if err != nil && !errors.As(err, &exitError) {
		return command.ProcessState.ExitCode(), fmt.Errorf("command wait: %w", err)
	}

	return command.ProcessState.ExitCode(), nil
}

func composeEnvSlice(env Environment) []string {
	envSlice := make([]string, 0, len(env))
	for key, value := range env {
		if !value.NeedRemove {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", key, value.Value))
		}
	}

	return envSlice
}
