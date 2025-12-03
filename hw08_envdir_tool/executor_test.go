package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd_EmptyCommand(t *testing.T) {
	var env Environment
	var out, errBuf bytes.Buffer

	_, err := RunCmd([]string{}, env, &errBuf, &out, nil, nil)
	assert.Error(t, err)
}

func TestRunCmd_EnvStdinStdoutStderr(t *testing.T) {
	env := Environment{
		"FOO":    {Value: "BAR"},
		"REMOVE": {Value: "XXX", NeedRemove: true},
	}

	var out, errBuf bytes.Buffer
	input := strings.NewReader("inval\n")

	cmd := []string{
		"sh", "-c",
		`read -r IN; printf %s "${IN}:$FOO"; >&2 printf err; printf :; printf %s "$REMOVE"`,
	}

	code, err := RunCmd(cmd, env, &errBuf, &out, input, nil)

	assert.NoError(t, err)
	assert.Equal(t, 0, code)
	assert.Equal(t, "inval:BAR:", out.String())
	assert.Equal(t, "err", errBuf.String())
}

func TestRunCmd_NonZeroExit_ReturnsError(t *testing.T) {

	var out, errBuf bytes.Buffer
	code, err := RunCmd([]string{"sh", "-c", "exit 7"}, Environment{}, &errBuf, &out, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 7, code)
}

//func TestRunCmd_SignalForwarding(t *testing.T) {
//	var out, errBuf bytes.Buffer
//	signals := make(chan syscall.Signal, 1)
//
//	cmd := []string{
//		"sh", "-c",
//		`trap "echo got signal; exit 5" INT TERM HUP; tail -f /dev/null & wait`,
//	}
//
//	done := make(chan struct{})
//	var code int
//	var runErr error
//
//	go func() {
//		code, runErr = RunCmd(cmd, Environment{}, &errBuf, &out, nil, signals)
//		close(done)
//	}()
//
//	// Даем процессу стартовать
//	time.Sleep(200 * time.Millisecond)
//	// Посылаем SIGINT
//	signals <- syscall.SIGINT
//
//	select {
//	case <-done:
//		// ok
//	case <-time.After(30 * time.Second):
//		t.Error("timeout waiting for command to exit after signal")
//	}
//
//	assert.NoError(t, runErr)
//	assert.Equal(t, 0, code)
//	assert.Equal(t, "trapped", out.String())
//	assert.Equal(t, 0, len(signals))
//}
