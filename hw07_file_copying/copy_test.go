package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		name     string
		offset   int
		limit    int
		expected string
	}{
		{
			"offset=0 limit=0",
			0,
			0,
			"testdata/out_offset0_limit0.txt",
		},
		{
			"offset=0 limit=10",
			0,
			10,
			"testdata/out_offset0_limit10.txt",
		},
		{
			"offset=0 limit=10",
			0,
			1000,
			"testdata/out_offset0_limit1000.txt",
		},
		{
			"offset=0 limit=10",
			0,
			10000,
			"testdata/out_offset0_limit10000.txt",
		},
		{
			"offset=0 limit=10",
			100,
			1000,
			"testdata/out_offset100_limit1000.txt",
		},
		{
			"offset=0 limit=10",
			6000,
			1000,
			"testdata/out_offset6000_limit1000.txt",
		},
	}

	fromFileName := "testdata/input.txt"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			toFileName, err := createTempFile(tc.limit, tc.offset)
			if err != nil {
				panic(err)
			}

			copyError := Copy(fromFileName, toFileName, int64(tc.offset), int64(tc.limit))

			expectedContent, exErr := os.ReadFile(tc.expected)
			if exErr != nil {
				panic(exErr)
			}

			actualContent, acErr := os.ReadFile(toFileName)
			if acErr != nil {
				panic(acErr)
			}

			assert.Nil(t, copyError, "function return error")
			assert.Equal(t, expectedContent, actualContent, "copied data is not like expected")
		})
	}

	t.Run("source file not exists", func(t *testing.T) {
		toFileName, err := createTempFile(-10, -10)
		if err != nil {
			panic(err)
		}
		copyErr := Copy("some not existed file name", toFileName, 0, 0)

		actualContent, acErr := os.ReadFile(toFileName)
		if acErr != nil {
			panic(acErr)
		}

		assert.Empty(t, actualContent)
		assert.NotNil(t, copyErr)
	})

	t.Run("offset is to much", func(t *testing.T) {
		toFileName, err := createTempFile(0, 100000)
		if err != nil {
			panic(err)
		}
		copyErr := Copy(fromFileName, toFileName, 100000, 0)

		actualContent, acErr := os.ReadFile(toFileName)
		if acErr != nil {
			panic(acErr)
		}

		assert.Empty(t, actualContent)
		assert.NotNil(t, copyErr)
	})

	t.Run("iregular file", func(t *testing.T) {
		toFileName, err := createTempFile(0, 100000)
		if err != nil {
			panic(err)
		}
		copyErr := Copy("/dev/random", toFileName, 0, 0)

		actualContent, acErr := os.ReadFile(toFileName)
		if acErr != nil {
			panic(acErr)
		}

		assert.Empty(t, actualContent)
		assert.NotNil(t, copyErr)
	})
}

func createTempFile(limit, offset int) (string, error) {
	f, err := os.CreateTemp("", fmt.Sprintf("temporary_result_limit%d_offset%d.txt", limit, offset))
	if err != nil {
		return "", fmt.Errorf("create temporary file: %w", err)
	}
	closeErr := f.Close()
	if closeErr != nil {
		return "", fmt.Errorf("close temporary file: %w", closeErr)
	}

	return f.Name(), nil
}
