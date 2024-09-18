package main

import (
	"fmt"
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

func handleErr(err *CmdError) {
	if err == nil { return }
	fmt.Fprintf(os.Stderr, "%s %s\n", Bold(Color("error:", RED)), err.Error())
	if err.Type == "new" || err.Type == "add" {
		fmt.Fprintf(os.Stderr, "\nrun %s for usage\n", Italic("gopher", Color("help", BLUE)))
	}
	os.Exit(1)
}

