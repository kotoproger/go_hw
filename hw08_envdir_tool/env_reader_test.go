package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDirSuccess(t *testing.T) {
	cases := []struct {
		name     string
		dir      string
		expected Environment
	}{
		{
			name: "simple die",
			dir:  "/testdata/env",
			expected: Environment{
				"BAR": EnvValue{
					Value:      "bar",
					NeedRemove: false,
				},
				"EMPTY": EnvValue{
					Value:      "",
					NeedRemove: true,
				},
				"FOO": EnvValue{
					Value:      "   foo\nwith new line",
					NeedRemove: false,
				},
				"HELLO": EnvValue{
					Value:      `"hello"`,
					NeedRemove: false,
				},
				"UNSET": EnvValue{
					Value:      ``,
					NeedRemove: true,
				},
			},
		},
		{
			name: "dir with subdirs",
			dir:  "/testdata/recurse",
			expected: Environment{
				"EMPTY": EnvValue{
					Value:      "",
					NeedRemove: true,
				},
				"FOO": EnvValue{
					Value:      "   foo\nwith new line",
					NeedRemove: false,
				},
				"HELLO": EnvValue{
					Value:      `"hello"`,
					NeedRemove: false,
				},
				"UNSET": EnvValue{
					Value:      ``,
					NeedRemove: true,
				},
			},
		},
		{
			name:     "empty directory",
			dir:      "/testdata/empty",
			expected: Environment{},
		},
	}

	for _, testCase := range cases {
		basePath, err := os.Getwd()
		os.Mkdir(basePath+"/testdata/empty", 0o755)
		if err != nil {
			panic(err)
		}
		t.Run(testCase.name, func(t *testing.T) {
			vars, err := ReadDir(basePath + testCase.dir)
			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, vars)
		})
	}
}

func TestReadDirFailure(t *testing.T) {
	basePath, err := os.Getwd()
	os.Chmod(basePath+"/testdata/failure/FOO", 0o055)
	defer os.Chmod(basePath+"/testdata/failure/FOO", 0o755)

	if err != nil {
		panic(err)
	}
	testCases := []struct {
		name string
		dir  string
	}{
		{
			name: "directory not exists",
			dir:  "/testdata/not_exist",
		},
		{
			name: "read file error",
			dir:  "/testdata/failure",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fmt.Println(basePath + testCase.dir)
			result, err := ReadDir(basePath + testCase.dir)

			var got *fs.PathError
			assert.True(t, errors.As(err, &got))
			var empty Environment
			assert.Equal(t, empty, result)
		})
	}
}
