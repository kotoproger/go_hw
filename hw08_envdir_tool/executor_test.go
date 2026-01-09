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

	_, err := RunCmd([]string{}, env, &errBuf, &out, nil)
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

	code, err := RunCmd(cmd, env, &errBuf, &out, input)

	assert.NoError(t, err)
	assert.Equal(t, 0, code)
	assert.Equal(t, "inval:BAR:", out.String())
	assert.Equal(t, "err", errBuf.String())
}

func TestRunCmd_NonZeroExit_ReturnsError(t *testing.T) {
	var out, errBuf bytes.Buffer
	code, err := RunCmd([]string{"sh", "-c", "exit 7"}, Environment{}, &errBuf, &out, nil)
	assert.NoError(t, err)
	assert.Equal(t, 7, code)
}
