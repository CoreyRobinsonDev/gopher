/*
	Main commands to be ran by gopher

new | add | test | build | run | help | config
*/
package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)


var pad string = "    "
// execute cmd
func RunCmd(cmd string, a... string) error {
	args := a[1:]
	switch cmd {
	case "new": return new(args[0])
	case "add": add()
	case "test": test()
	case "build": return build()
	case "run": run()
	case "help": help()
	case "config": help()
	default: return errors.New(fmt.Sprintf("no such command: %s", cmd))
	}

	return nil
}

// create new go module
func new(path string) error {
	if !strings.Contains(path, "/") {
		return errors.New("go module requires repository location\nexample: gopher new github.com/user/mymodule")
	}
	pathArr := strings.Split(path, "/")
	name := pathArr[len(pathArr) - 1]

	unwrap(fmt.Printf("%sCreating binary `%s` module\n", pad, name))
	expect(os.Mkdir(name, 0755))
	expect(os.Chdir(name))
	goCmd := exec.Command("go", "mod", "init", path)
	gitCmd := exec.Command("git", "init")
	expect(goCmd.Run())
	expect(gitCmd.Run())
	expect(os.Mkdir("src", 0755))
	expect(os.Mkdir("target", 0755))

	f1 := unwrap(os.Create("./src/main.go"))
	f2 := unwrap(os.Create("./.gitignore"))
	defer f1.Close()
	defer f2.Close()
	unwrap(f1.WriteString("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}"))
	unwrap(f2.WriteString("/target"))

	return nil
}

// build go binaries for mac, linux, and windows
func build() error {
	envCmd := exec.Command("go", "env")
	grepOSCmd := exec.Command("grep", "GOHOSTOS")
	grepARCHCmd := exec.Command("grep", "GOHOSTARCH")
	out := unwrap(envCmd.Output())
	grepOSCmdIn := unwrap(grepOSCmd.StdinPipe())
	grepARCHCmdIn := unwrap(grepARCHCmd.StdinPipe())
	grepOSCmdOut := unwrap(grepOSCmd.StdoutPipe())
	grepARCHCmdOut := unwrap(grepARCHCmd.StdoutPipe())
	expect(grepOSCmd.Start())
	expect(grepARCHCmd.Start())
	grepOSCmdIn.Write(out)
	grepARCHCmdIn.Write(out)
	grepARCHCmdIn.Close()
	grepOSCmdIn.Close()
	opsysList := [][]string {
		{"windows", "amd64"},
		{"windows", "arm64"},
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
	}
	
	OSout := string(unwrap(io.ReadAll(grepOSCmdOut)))
	ARCHout := string(unwrap(io.ReadAll(grepARCHCmdOut)))
	grepARCHCmd.Wait()
	grepOSCmd.Wait()

	osv := strings.Trim(strings.Split(OSout, "=")[1], " '\n\t")
	arch := strings.Trim(strings.Split(ARCHout, "=")[1], " '\n\t")
	
	for _, item := range opsysList {
		sysop := item[0]
		sysarch := item[1]
		if sysop == osv && sysarch == arch { continue }
		os.Setenv("GOOS", sysop)
		os.Setenv("GOARCH", sysarch)
		name := fmt.Sprintf()
		buildCmd := exec.Command("go", "build", "-o")
	}

	return nil
}

func add() {}
func test() {}
func run() {}
func help() {}
func config() {}

func unwrap[T any](val T, err error) T {
	if err != nil { panic(err) }

	return val
}

func unwrapOr[T any](val T, err error) func(T) T {
	if err != nil {
		return func(d T) T {
			return d
		}
	} else {
		return func(_ T) T {
			return val
		}
	}
}

func expect(err error) {
	if err != nil { panic(err) }
}
