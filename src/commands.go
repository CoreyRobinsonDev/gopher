package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)



func RunCmd(cmd string, a... string) *CmdError {
	var args []string
	if len(a) > 0 {
		args = a[1:]
	}
	switch cmd {
	case "new": 
		if len(args) == 0 {
			return &CmdError {
				Type: CmdNew,
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
				Type: CmdAdd,
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
	case "version": version()
	default: return &CmdError {
		Type: CmdInvalid,
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
	Unwrap(fmt.Printf("%sCreated binary '%s' module\n", PAD, name))
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
	var oparchPairs [][]string 
	var announceBuild bool
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
	var datArr []string
	if strings.Contains(dat, "/") {
		datArr = strings.Split(dat, "/")
	} else {
		datArr = strings.Split(dat, " ")
	}
	module := datArr[len(datArr) - 1]
	
	OSout := string(Unwrap(io.ReadAll(grepOSCmdOut)))
	ARCHout := string(Unwrap(io.ReadAll(grepARCHCmdOut)))
	grepARCHCmd.Wait()
	grepOSCmd.Wait()

	osv := strings.Trim(strings.Split(OSout, "=")[1], " '\n\t")
	arch := strings.Trim(strings.Split(ARCHout, "=")[1], " '\n\t")
	
	if len(args) > 0 {
		if args[0] == "--cross-platform" || args[0] == "-x" {
			oparchPairs = Unwrap(GetPreference[[][]string](PrefOpArchPairs))
			announceBuild = true
		} else {
			oparchPairs = [][]string{{osv, arch}}
		}
	} else {
		oparchPairs = [][]string{{osv, arch}}
	}
	var wg sync.WaitGroup
	errs := make(chan error, len(oparchPairs))
	mainErr := make(chan string, 1)

	for _, item := range oparchPairs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var name string
			var buildCmd *exec.Cmd
			sysop := item[0]
			sysarch := item[1]
			if announceBuild {
				fmt.Printf(
					"%sBuilding %s binary for %s architecture\n", 
					PAD, sysop, sysarch,
				)
			}

			if osv == sysop && arch == sysarch {
				name = fmt.Sprintf(
					"./bin/%s",
					module,
				)
				if len(args) > 0 && !slices.Contains(args, "-x") && !slices.Contains(args, "--cross-platform") {
					args = append([]string{"build"}, args...)
					buildCmd = exec.Command("go", args...)
				} else {
					buildCmd = exec.Command("go", "build", "-o", name, "./src/")
				}

				o, _ := buildCmd.CombinedOutput()

				if len(o) > 1 {
					mainErr <- string(o)
				}
			} else {
				name = fmt.Sprintf(
					"./bin/%s-%s-%s",
					module,
					sysarch,
					sysop,	
				)
				buildCmdStr := fmt.Sprintf(
					"GOOS=%s GOARCH=%s go build -o %s ./src/",
					sysop,
					sysarch,
					name, 
				)
				buildCmd = exec.Command("bash", "-c", buildCmdStr)

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

	if len(errStr) > 0 || len(mainErr) > 0 {
		if len(mainErr) > 0 {
			fmt.Print(<-mainErr)
		} else {
			return &CmdError {
				Type: CmdBuild,
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
			fmt.Printf("Add dependencies to current module and install them.\n\nWhen a full package name isn't provided %s will do a search on pkg.go.dev for matching packages. The number of results returned on this search can be adjusted with %s.\n\n%s %s\n%s %s\n",
				Bold("gopher add"),
				Bold("gopher config"),
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("add", BLUE),
					"rsc.io/quote",
				),
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("add", BLUE),
					"gofiber",
				),
			)
		case "test": 
		case "build": 
			fmt.Printf("compile packages and dependencies\n\n%s should be executed at the root of your module and will expect the entry point of your program to be ./src/main.go\n\n%s %s\n\n%s\n%s\n\n%s %s\n",
				Bold("gopher build"),
				Bold(Color("usage:", PURPLE)),
				Italic(
					"gopher",
					Color("build", BLUE),
					Color("[...ARGS]", CYAN),
				),
				Bold(Color("arguments:", PURPLE)),
				PAD + "-x,--cross-platform" + "\t\t" + "build binaries for seperate operating systems and cpu architectures speficied by your gopher configuration",
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("build", BLUE),
				),
			)
		case "run": 
			fmt.Printf("compile and run Go program\n\n%s should be executed at the root of your module and will expect the entry point of your program to be ./src/main.go\n\n%s %s\n\n%s %s\n",
				Bold("gopher run"),
				Bold(Color("usage:", PURPLE)),
				Italic(
					"gopher",
					Color("run", BLUE),
					Color("[...ARGS]", CYAN),
				),
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
			Type: CmdHelp,
			Msg: fmt.Sprintf("no such command: %s", cmd),
		}
		}
	} else {
		fmt.Printf("A Go module manager\n\n%s %s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n\nsee %s for more information about a specific command\n",
			Bold(Color("usage:", PURPLE)),
			Italic(
				"gopher",
				Color("[COMMAND]", BLUE),
				Color("[...ARGS]", CYAN),
			),
			Bold(Color("commands:", PURPLE)),
			PAD + "add" + "\t\t" + "add dependencies to current module and install them",
			PAD + "build" + "\t" + "compile packages and dependencies",
			PAD + "config" + "\t" + "configure gopher settings",
			PAD + "help" + "\t" + "this",
			PAD + "new" + "\t\t" + "create new go module",
			PAD + "run" + "\t\t" + "compile and run Go program",
			PAD + "test" + "\t" + "run Go test packages",
			PAD + "version" + "\t" + "print Go version",
			Italic(
				"gopher",
				Color("help", BLUE),
				Color("[COMMAND]", CYAN),
			),
		)
	}

	return nil
}

func add(pkg string) *CmdError {
	if strings.Contains(pkg, "/") {
		getCmd := exec.Command("go", "get", pkg)
		out, _ := getCmd.CombinedOutput()
		fmt.Print(string(out))
	} else {
		horizonalCharLimit := 80
		pkgQueryLimit := Unwrap(GetPreference[int](PrefPkgQueryLimit))
		if pkgQueryLimit > 100 { 
			return &CmdError {
				Type: CmdAdd,
				Msg: fmt.Sprintf("PkgQueryLimit (%d) >100", pkgQueryLimit),
			}
		} 
		url := fmt.Sprintf(
			"https://pkg.go.dev/search?limit=%d&m=package&q=%s",
			pkgQueryLimit,
			pkg,
		)
		res := Unwrap(http.Get(url))
		dat := string(Unwrap(io.ReadAll(res.Body)))
		defer res.Body.Close()
		lines := strings.Split(dat, "\n")
		pkgHTML := [][]string{}
		pkgNames := []string{}
		lineNums := []int{}

		for i, line := range lines {
			if strings.Contains(line, "\"SearchSnippet\"") {
				lineNums = append(lineNums, i)
			}
			if i == len(lines) - 1 {
				lineNums = append(lineNums, i)
			}
		}

		for i := range lineNums {
			if i == len(lineNums) - 1 {break}
			pkgHTML = append(pkgHTML, lines[lineNums[i]:lineNums[i+1]])
		}

		for i := len(pkgHTML)-1; i >= 0; i-- {
			var pkgName string
			var pkgMeta string
			var pkgDesc string
			for ii, pkgLine := range pkgHTML[i] {
				pkgLine = strings.Trim(pkgLine, " \t\n")
				if strings.Contains(pkgLine, "SearchSnippet-header-path") {
					write := false
					for _, ch := range pkgLine {
						if ch == '(' && write == false { 
							write = true; continue 
						} else if ch == ')' && write == true { write = false }
						if write {
							pkgName += string(ch)
						}
					}
				} else if strings.Contains(pkgLine, "published on") {
					write := false
					for _, ch := range pkgLine {
						if ch == '>' && write == false { 
							write = true; continue 
						} else if ch == '<' && write == true { write = false }
						if write {
							pkgMeta += string(ch)
						}
					}
				} else if strings.Contains(pkgLine, "SearchSnippet-synopsis") {
					pkgDesc = strings.Trim(pkgHTML[i][ii+1], " \t\n")
				}
			}
			pkgMetaArr := strings.Split(pkgMeta, " ")
			version := pkgMetaArr[0]
			fmt.Printf(
				"%s %s %s %s\n", 
				Color(strconv.Itoa(i+1), BLUE),
				Bold(pkgName),
				Color(version, CYAN),
				Bold("("+strings.Join(pkgMetaArr[1:], " ")+")"),
			)
			pkgDesc = strings.ReplaceAll(pkgDesc, "&#34;", "\"")
			pkgDesc = strings.ReplaceAll(pkgDesc, "&#39;", "'")
			if len(pkgDesc) > 0 {
				substring := ""
				for _, ch := range pkgDesc {
					substring += string(ch)
					if len(substring) == horizonalCharLimit {
						substringArr := strings.Split(substring, " ")
						fmt.Printf("%s%s\n", PAD, 
							strings.Join(substringArr[:len(substringArr)-1], " "),
						)
						substring = substringArr[len(substringArr)-1]
					}
				}
				if len(substring) > 0 {
					fmt.Printf("%s%s\n", PAD, substring)
				}
			}
			pkgNames = append([]string{pkgName}, pkgNames...)
		}
		fmt.Printf(
			"%s Packages to install (eg: 1 2 3)\n",
			Bold(Color("==>", YELLOW)),
		)
		fmt.Printf(
			"%s",
			Bold(Color("==> ", YELLOW)),
		)
		reader := bufio.NewReader(os.Stdin)
		in := Unwrap(reader.ReadString('\n'))
		in = in[:len(in)-1]
		opt, err := strconv.Atoi(in)
		if err != nil || opt >= pkgQueryLimit {
			return &CmdError {
				Type: CmdAdd,
				Msg: fmt.Sprintf("index '%s' not found\n\nenter a numeric value from 1-%d",
					in,
					pkgQueryLimit-1,
				),
			}
		}

		output := make(chan string)
		go func() {
			getCmd := exec.Command("go", "get", pkgNames[opt-1])
			o, _ := getCmd.CombinedOutput()
			output <- string(o)
		}()
		count := 0
		loop: for {
			count++
			fmt.Printf(
				"downloading %s %s\r",
				Bold(pkgNames[opt-1]),
				Color("("+strconv.Itoa(count)+"s)", BLUE),
			)
			time.Sleep(1 * time.Second)
			select {
			case out := <-output:
				fmt.Printf("\n%s",out)
				close(output)
				break loop
			default:
			}
		}
	}

	return nil
}

func test() {}
func tidy() {}
func version() {}

var PAD string = "    "
var DEFAULT_PREFERENCES = `PkgQueryLimit=10
OpArchPairs=windows,amd64,windows,arm64,linux,amd64,linux,arm64,darwin,amd64,darwin,arm64`

type Command int
const (
	CmdNew = iota
	CmdAdd
	CmdHelp
	CmdTidy
	CmdBuild
	CmdConfig
	CmdRun
	CmdVersion
	CmdInvalid
)

var commandName = map[Command]string {
	CmdNew: "new",
	CmdAdd: "add",
	CmdHelp: "help",
	CmdTidy: "tidy",
	CmdBuild: "build",
	CmdConfig: "config",
	CmdRun: "run",
	CmdVersion: "version",
	CmdInvalid: "invalid",
}
func (c Command) String() string {
	return commandName[c]
}

type CmdError struct {
	Type Command
	Msg string
}

func (c CmdError) Error() string {
	return c.Msg
}

