package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/CoreyRobinsonDev/gopher/src/commands"
)


func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		handleErr(errors.New("no arguments"))
	}

	err := commands.RunCmd(args[0], args...)
	handleErr(err)
}

func handleErr(err error) {
	if err == nil { return }
	fmt.Fprintln(os.Stderr, "error:", err.Error())
	fmt.Fprintln(os.Stderr, "\nrun 'gopher help' for usage")
	os.Exit(1)
}

