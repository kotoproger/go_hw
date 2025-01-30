package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")
	Digits           = map[int]bool{
		0: true,
		1: true,
		2: true,
		3: true,
		4: true,
		5: true,
		6: true,
		7: true,
		8: true,
		9: true,
	}
)

func Unpack(input string) (string, error) {
	var lastSymbol string
	var result strings.Builder
	var escaped bool
	for _, RuneSymbol := range input {
		intValue, err := strconv.Atoi(string(RuneSymbol))
		if err == nil && Digits[intValue] && !escaped {
			if lastSymbol == "" {
				return "", ErrInvalidString
			}
			result.WriteString(repeatRune(lastSymbol, intValue))
			lastSymbol = ""
		} else {
			if !escaped && string(RuneSymbol) == "\\" {
				escaped = true

				continue
			}
			if escaped && err != nil && string(RuneSymbol) != "\\" {
				return "", ErrInvalidString
			}
			escaped = false
			result.WriteString(lastSymbol)
			lastSymbol = string(RuneSymbol)
		}
	}
	result.WriteString(lastSymbol)

	return result.String(), nil
}

func repeatRune(symbol string, counter int) string {
	return strings.Repeat(symbol, counter)
}
