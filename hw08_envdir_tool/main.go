package main

import (
	"fmt"
	"os"
)

func main() {
	envs, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	code, runErr := RunCmd(
		os.Args[2:],
		envs,
		os.Stderr,
		os.Stdout,
		os.Stdin,
	)
	if runErr != nil {
		fmt.Println(runErr)
		os.Exit(runErr.Code())
	}
	os.Exit(code)
}
