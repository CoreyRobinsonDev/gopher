package main

import (
	"os"
)


func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		err := RunCmd("help")
		handleErr(err)
		os.Exit(0)
	}

	err := RunCmd(args[0], args...)
	handleErr(err)
}

