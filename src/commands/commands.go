package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)


var pad string = "    "
var opsysList [][]string = [][]string {
	{"windows", "amd64"},
	{"windows", "arm64"},
	{"linux", "amd64"},
	{"linux", "arm64"},
	{"darwin", "amd64"},
	{"darwin", "arm64"},
}

type CmdError struct {
	Type string
	Msg string
}

func (c CmdError) Error() string {
	return c.Msg
}

// execute cmd
func RunCmd(cmd string, a... string) *CmdError {
	var args []string
	if len(a) > 0 {
		args = a[1:]
	}
	switch cmd {
	case "new": 
		if len(args) == 0 {
			return &CmdError {
				Type: "new",
				Msg: "go module requires repository location\n\nexample: gopher new github.com/user/mymodule",
			}
		}
		return new(args[0])
	case "add": 
		if len(args) == 0 {
			return &CmdError {
				Type: "add",
				Msg: "no go module specified\n\nexample: gopher add rsc.io/quote",
			}
		}
		return add(args[0])
	case "test": test()
	case "build": return build()
	case "run": return run()
	case "help": 
		if len(args) == 0 {
			return help("")
		}
		return help(args[0], args[1:]...)
	case "config": config()
	default: return &CmdError {
		Type: "",
		Msg: fmt.Sprintf("no such command: %s", cmd),
	}
	}

	return nil
}

// create new go module
func new(path string) *CmdError {
	// if !strings.Contains(path, "/") {
	// 	return &CmdError {
	// 		Type: "new",
	// 		Msg: "go module requires repository location\nexample: gopher new github.com/user/mymodule",
	// 	}
	// }
	// pathArr := strings.Split(path, "/")
	// name := pathArr[len(pathArr) - 1]

	name := path
	Unwrap(fmt.Printf("%sCreating binary `%s` module\n", pad, name))
	Expect(os.Mkdir(name, 0755))
	Expect(os.Chdir(name))
	goCmd := exec.Command("go", "mod", "init", path)
	gitCmd := exec.Command("git", "init")
	Expect(goCmd.Run())
	Expect(gitCmd.Run())
	Expect(os.Mkdir("src", 0755))
	Expect(os.Mkdir("bin", 0755))

	f1 := Unwrap(os.Create("./src/main.go"))
	f2 := Unwrap(os.Create("./.gitignore"))
	f3 := Unwrap(os.Create("./README.md"))
	defer f1.Close()
	defer f2.Close()
	defer f3.Close()
	Unwrap(f1.WriteString("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}"))
	Unwrap(f2.WriteString("/bin"))
	Unwrap(f3.WriteString(fmt.Sprintf("# %s", name)))

	return nil
}

// build go binaries for mac, linux, and windows
func build() *CmdError {
	envCmd := exec.Command("go", "env")
	grepOSCmd := exec.Command("grep", "GOHOSTOS")
	grepARCHCmd := exec.Command("grep", "GOHOSTARCH")
	out := Unwrap(envCmd.Output())
	grepOSCmdIn := Unwrap(grepOSCmd.StdinPipe())
	grepARCHCmdIn := Unwrap(grepARCHCmd.StdinPipe())
	grepOSCmdOut := Unwrap(grepOSCmd.StdoutPipe())
	grepARCHCmdOut := Unwrap(grepARCHCmd.StdoutPipe())
	Expect(grepOSCmd.Start())
	Expect(grepARCHCmd.Start())
	grepOSCmdIn.Write(out)
	grepARCHCmdIn.Write(out)
	grepARCHCmdIn.Close()
	grepOSCmdIn.Close()

	dat := string(Unwrap(os.ReadFile("./go.mod")))
	dat = strings.Split(dat, "\n")[0]
	datArr := strings.Split(dat, "/")
	module := datArr[len(datArr) - 1]
	
	OSout := string(Unwrap(io.ReadAll(grepOSCmdOut)))
	ARCHout := string(Unwrap(io.ReadAll(grepARCHCmdOut)))
	grepARCHCmd.Wait()
	grepOSCmd.Wait()

	osv := strings.Trim(strings.Split(OSout, "=")[1], " '\n\t")
	arch := strings.Trim(strings.Split(ARCHout, "=")[1], " '\n\t")
	
	var wg sync.WaitGroup
	errs := make(chan error, len(opsysList))
	mainErr := make(chan string, 1)

	for _, item := range opsysList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var name string
			sysop := item[0]
			sysarch := item[1]
			fmt.Printf(
				"%sBuilding %s binary for %s architecture\n", 
				pad, sysop, sysarch,
			)

			if osv == sysop && arch == sysarch {
				name = fmt.Sprintf(
					"./bin/%s",
					module,
				)
			} else {
				name = fmt.Sprintf(
					"./bin/%s-%s-%s",
					module,
					sysarch,
					sysop,	
				)
			}
			buildCmdStr := fmt.Sprintf(
				"GOOS=%s GOARCH=%s go build -o %s ./src/main.go",
				sysop,
				sysarch,
				name, 
			)
			buildCmd := exec.Command("bash", "-c", buildCmdStr)

			if osv == sysop && arch == sysarch {
				o, _ := buildCmd.CombinedOutput()

				if len(o) > 0 {
					mainErr <- fmt.Sprintf("-----------------------------------------------------\n%s", string(o))
				}
			} else {
				if err := buildCmd.Run(); err != nil {
					errs <- errors.New(
						fmt.Sprintf("could not build %s binary for %s architecture",
							sysop, sysarch,
						),
					)
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	close(mainErr)

	var errStr string
	for err := range errs {
		errStr += err.Error() + "\n"
	}

	if len(errStr) > 0 {
		if len(mainErr) > 0 {
			return &CmdError {
				Type: "build",
				Msg: errStr + <-mainErr,
			}
		} else {
			return &CmdError {
				Type: "build",
				Msg: errStr[:len(errStr)-1],
			}
		}
	}

	return nil
}

func run() *CmdError {
	runCmd := exec.Command("go", "run", "./src/")
	o, _ := runCmd.CombinedOutput()
	fmt.Print(string(o))

	return nil
}

func help(cmd string, moreCmds ...string) *CmdError {
	if len(cmd) > 0 {
		switch cmd {
		case "new": 
			fmt.Println(`create new go module

example: gopher new github.com/user/mymodule`)
		case "add": 
		case "test": 
		case "build": 
			fmt.Println(`compile packages and dependencies

'gopher build' should be executed at the root of your module and will expect the entry point of your program to be ./src/main.go

example: gopher build`)
		case "run": 
			fmt.Println(`compile and run Go program

'gopher run' should be executed at the root of your module and will expect the entry point of your program to be ./src/main.go

example: gopher run`)
		case "help": 
			if len(moreCmds) > 0 {
				help(moreCmds[0], moreCmds[1:]...)
			} else {
				help("")
			}
		case "config": 
		default: return &CmdError {
			Type: "help",
			Msg: fmt.Sprintf("no such command: %s", cmd),
		}
		}
	} else {
		fmt.Println(`A Go module manager

usage: gopher [COMMAND] [...ARGS]

commands:
    add	        add dependencies to current module and install them
    build       compile packages and dependencies
    config      configure gopher settings
    help        this
    new	        create new go module
    run	        compile and run Go program
    test        run Go test packages
    version     print Go version

see 'gopher help [COMMAND]' for more information about a specific command`)
	}

	return nil
}

func add(pkg string) *CmdError {
	getCmd := exec.Command("go", "get", pkg)
	out, _ := getCmd.CombinedOutput()
	fmt.Print(string(out))

	return nil
}

func test() {}
func config() {}

