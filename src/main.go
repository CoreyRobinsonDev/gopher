package main

import (
	"fmt"
	"os"

	"github.com/CoreyRobinsonDev/gopher/src/commands"
)


func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		err := commands.RunCmd("help")
		handleErr(err)
		os.Exit(0)
	}

	err := commands.RunCmd(args[0], args...)
	handleErr(err)
}

func handleErr(err *commands.CmdError) {
	if err == nil { return }
	fmt.Fprintf(os.Stderr, "%s %s\n", commands.Bold(commands.Color("error:", commands.RED)), err.Error())
	if err.Type == "new" || err.Type == "add" {
		fmt.Fprintf(os.Stderr, "\nrun %s for usage\n", commands.Italic("gopher", commands.Color("help", commands.BLUE)))
	}
	os.Exit(1)
}

