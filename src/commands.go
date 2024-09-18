package main

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
				Msg: fmt.Sprintf("go module requires repository location\n\n%s %s",
					Bold(Color("example:", PURPLE)),
					Italic(
						"gopher",
						Color("new", BLUE),
						"github.com/user/mymodule",
					),
				),
			}
		}
		return new(args[0])
	case "add": 
		if len(args) == 0 {
			return &CmdError {
				Type: "add",
				Msg: fmt.Sprintf("no go module specified\n\n%s %s",
					Bold(Color("example:", PURPLE)),
					Italic(
						"gopher",
						Color("add", BLUE),
						"rsc.io/quote",
					),
				),
			}
		}
		return add(args[0])
	case "test": test()
	case "build": return build(args...)
	case "run": return run(args...)
	case "help": 
		if len(args) == 0 {
			return help("")
		}
		return help(args[0], args[1:]...)
	case "config": config()
	case "version": version()
	default: return &CmdError {
		Type: "",
		Msg: fmt.Sprintf("no such command: %s", cmd),
	}
	}

	return nil
}

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
	Unwrap(f2.WriteString("bin/\n*.exe\n*.exe~\n*.dll\n*.so\n*.dylib\n\n*.test\n*.out\nvendor/\n\ngo.work\ngo.work.sum\n\n.env"))
	Unwrap(f3.WriteString(fmt.Sprintf("# %s", name)))

	return nil
}

func build(args ...string) *CmdError {
	var opsysPairs [][]string 
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
	
	if len(args) > 0 {
		if args[0] == "--cross-platform" || args[0] == "-x" {
			opsysPairs = opsysList
		}
	} else {
		opsysPairs = [][]string{{osv, arch}}
	}
	var wg sync.WaitGroup
	errs := make(chan error, len(opsysPairs))
	mainErr := make(chan string, 1)

	for _, item := range opsysPairs {
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
				"GOOS=%s GOARCH=%s go build -o %s ./src/",
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

func run(args ...string) *CmdError {
	args = append([]string {"run", "./src/"}, args...)
	runCmd := exec.Command("go", args...)
	out, _ := runCmd.CombinedOutput()
	fmt.Print(string(out))

	return nil
}

func help(cmd string, moreCmds ...string) *CmdError {
	if len(cmd) > 0 {
		switch cmd {
		case "new": 
			fmt.Printf("create new go module\n\n%s %s\n",
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("new", BLUE),
					"github.com/user/mymodule",
				),
			)
		case "add": 
			fmt.Printf("add dependencies to current module and install them\n\n%s %s\n",
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("add", BLUE),
					"rsc.io/quote",
				),
			)
		case "test": 
		case "build": 
			fmt.Printf("compile packages and dependencies\n\n%s %s\n\n%s %s\n",
				Bold("gopher build"),
				"should be executed at the root of your module and will expect the entry point of your program to be ./src/main.go",
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("build", BLUE),
				),
			)
		case "run": 
			fmt.Printf("compile and run Go program\n\n%s %s\n\n%s %s\n",
				Bold("gopher run"),
				"should be executed at the root of your module and will expect the entry point of your program to be ./src/main.go",
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("run", BLUE),
				),
			)
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
		fmt.Printf("A Go module manager\n\n%s %s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n\nsee %s for more information about a specific command",
			Bold(Color("usage:", PURPLE)),
			Italic(
				"gopher",
				Color("[COMMAND]", BLUE),
				Color("[...ARGS]", CYAN),
			),
			Bold(Color("commands:", PURPLE)),
			pad + "add" + "\t\t" + "add dependencies to current module and install them",
			pad + "build" + "\t" + "compile packages and dependencies",
			pad + "config" + "\t" + "configure gopher settings",
			pad + "help" + "\t" + "this",
			pad + "new" + "\t\t" + "create new go module",
			pad + "run" + "\t\t" + "compile and run Go program",
			pad + "test" + "\t" + "run Go test packages",
			pad + "version" + "\t" + "print Go version",
			Italic(
				"gopher",
				Color("help", BLUE),
				Color("[COMMANDS]", CYAN),
			),
		)
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
func version() {}

