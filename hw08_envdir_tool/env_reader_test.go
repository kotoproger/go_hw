package main

import (
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
					Value:      " ",
					NeedRemove: false,
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
					Value:      " ",
					NeedRemove: false,
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
		os.Mkdir(basePath+"/testdata/empty", 0755)
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
