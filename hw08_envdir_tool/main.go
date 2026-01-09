package main

import "os"

func main() {
	envs, err := ReadDir(os.Args[1])
	if err != nil {
		panic(err)
	}
	code, err := RunCmd(
		os.Args[2:],
		envs,
		os.Stderr,
		os.Stdout,
		os.Stdin,
	)
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}
