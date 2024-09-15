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
	f3 := unwrap(os.Create("./README.md"))
	defer f1.Close()
	defer f2.Close()
	defer f3.Close()
	unwrap(f1.WriteString("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}"))
	unwrap(f2.WriteString("/target"))
	dat := string(unwrap(os.ReadFile("./go.mod")))
	dat = strings.Split(dat, "\n")[0]
	datArr := strings.Split(dat, "/")
	module := datArr[len(datArr) - 1]
	unwrap(f3.WriteString(fmt.Sprintf("# %s", module)))

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

	dat := string(unwrap(os.ReadFile("./go.mod")))
	dat = strings.Split(dat, "\n")[0]
	datArr := strings.Split(dat, "/")
	module := datArr[len(datArr) - 1]
	
	OSout := string(unwrap(io.ReadAll(grepOSCmdOut)))
	ARCHout := string(unwrap(io.ReadAll(grepARCHCmdOut)))
	grepARCHCmd.Wait()
	grepOSCmd.Wait()

	osv := strings.Trim(strings.Split(OSout, "=")[1], " '\n\t")
	arch := strings.Trim(strings.Split(ARCHout, "=")[1], " '\n\t")
	
	var wg sync.WaitGroup
	errs := make(chan error, len(opsysList))

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
					"./target/%s",
					module,
				)
			} else {
				name = fmt.Sprintf(
					"./target/%s-%s-%s",
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
			if err := buildCmd.Run(); err != nil {
				errs <- errors.New(
					fmt.Sprintf("could not build %s binary for %s architecture",
						sysop, sysarch,
					),
				)
			}
		}()
	}
	wg.Wait()
	close(errs)

	var errStr string
	for err := range errs {
		errStr += err.Error() + "\n"
	}

	if len(errStr) > 0 {
		return errors.New(errStr[:len(errStr)-1])
	}

	return nil
}

func run() {}
func add() {}
func test() {}
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
