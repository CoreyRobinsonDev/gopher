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

func RunCmd(cmd string, a ...string) *CmdError {
	var args []string
	if len(a) > 0 {
		args = a[1:]
	}
	switch cmd {
	case "new":
		if len(args) == 0 {
			return &CmdError{
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
			return &CmdError{
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
	case "test":
		return test(args...)
	case "build":
		return build(args...)
	case "run":
		return run(args...)
	case "tidy":
		return tidy()
	case "help":
		if len(args) == 0 {
			return help("")
		}
		return help(args[0], args[1:]...)
	case "version":
		return version()
	default:
		return &CmdError{
			Type: CmdInvalid,
			Msg:  fmt.Sprintf("no such command: %s", cmd),
		}
	}
}

func new(path string) *CmdError {
	name := ""
	if strings.Contains(path, "/") {
		pathArr := strings.Split(path, "/")
		name = pathArr[len(pathArr)-1]
	} else {
		name = path
	}

	Unwrap(fmt.Printf("%sCreated binary '%s' module\n", PAD, name))
	Expect(os.Mkdir(name, 0755))
	Expect(os.Chdir(name))
	goCmd := exec.Command("go", "mod", "init", path)
	gitCmd := exec.Command("git", "init")
	Expect(goCmd.Run())
	Expect(gitCmd.Run())

	f1 := Unwrap(os.Create("./main.go"))
	f2 := Unwrap(os.Create("./.gitignore"))
	f3 := Unwrap(os.Create("./README.md"))
	defer f1.Close()
	defer f2.Close()
	defer f3.Close()
	Unwrap(f1.WriteString("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"))
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
	module := datArr[len(datArr)-1]

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
					buildCmd = exec.Command("go", "build", "-o", name, ".")
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
					"GOOS=%s GOARCH=%s go build -o %s .",
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
			return &CmdError{
				Type: CmdBuild,
				Msg:  errStr[:len(errStr)-1],
			}
		}
	}

	return nil
}

func run(args ...string) *CmdError {
	args = append([]string{"run", "."}, args...)
	runCmd := exec.Command("go", args...)
	output, e := runCmd.CombinedOutput()

	if Unwrap(GetPreference[bool](PrefPrettyPrint)) && e != nil {
		outputLines := strings.Split(string(output), "\n")
		for _, line := range outputLines {
			if !strings.Contains(line, ":") {
				continue
			}
			arr := strings.Split(line, ":")
			file := arr[0]
			rownum := Unwrap(strconv.Atoi(arr[1]))
			colnum := Unwrap(strconv.Atoi(arr[2]))
			err := strings.Trim(strings.Join(arr[3:], ":"), " \t\n")
			f := Unwrap(os.Open("./" + file))
			previewLines := Unwrap(GetPreference[int](PrefPrettyPrintPreviewLines))
			defer f.Close()
			reader := bufio.NewScanner(f)

			fmt.Printf("[%s]\n", Color(file, BLUE))
			linenum := 0
			for reader.Scan() {
				linenum++
				if dif := linenum - rownum; dif <= previewLines && dif >= -previewLines {
					tabs := strings.Count(reader.Text(), "\t")
					if linenum == rownum {
						fmt.Printf("%s %s%s\n",
							Color(strconv.Itoa(linenum), GRAY),
							strings.Repeat("  ", tabs),
							Italic(strings.Trim(reader.Text(), " \t")),
						)
					} else {
						fmt.Printf("%s %s%s\n",
							Color(strconv.Itoa(linenum), GRAY),
							strings.Repeat("  ", tabs),
							strings.Trim(reader.Text(), " \t"),
						)
					}
				}
				if linenum == rownum {
					tabs := strings.Count(reader.Text(), "\t")
					pad := strings.Repeat(" ", colnum+len(strconv.Itoa(linenum))-tabs)
					fmt.Printf("%s%s%s\n", strings.Repeat("  ", tabs), pad, Color("^ "+err, RED))
				}
				if linenum > rownum+previewLines {
					fmt.Println()
					break
				}
			}
		}
	} else {
		fmt.Print(string(output))
	}

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
			fmt.Printf("run _test.go files\n\n%s %s\n",
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("test", BLUE),
				),
			)
		case "build":
			fmt.Printf("compile packages and dependencies\n\n%s should be executed at the root of your module and will expect the entry point of your program to be main.go\n\n%s %s\n\n%s\n%s\n\n%s %s\n",
				Bold("gopher build"),
				Bold(Color("usage:", PURPLE)),
				Italic(
					"gopher",
					Color("build", BLUE),
					Color("[...ARGS]", CYAN),
				),
				Bold(Color("arguments:", PURPLE)),
				PAD+"-x,--cross-platform"+"\t\t"+"build binaries for separate operating systems and cpu architectures speficied by your gopher configuration",
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("build", BLUE),
				),
			)
		case "run":
			fmt.Printf("compile and run Go program\n\n%s should be executed at the root of your module and will expect the entry point of your program to be main.go\n\n%s %s\n\n%s %s\n",
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
		case "version":
			fmt.Printf("print Go version\n\n%s %s\n",
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("version", BLUE),
				),
			)
		case "tidy":
			fmt.Printf("add missing and remove unused modules\n\n%s %s\n",
				Bold(Color("example:", PURPLE)),
				Italic(
					"gopher",
					Color("tidy", BLUE),
				),
			)
		default:
			return &CmdError{
				Type: CmdHelp,
				Msg:  fmt.Sprintf("no such command: %s", cmd),
			}
		}
	} else {
		fmt.Printf("A Go project manager\n\n%s %s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n\nsee %s for more information about a specific command\n",
			Bold(Color("usage:", PURPLE)),
			Italic(
				"gopher",
				Color("[COMMAND]", BLUE),
				Color("[...ARGS]", CYAN),
			),
			Bold(Color("commands:", PURPLE)),
			PAD+"add"+"\t\t"+"add dependencies to current module and install them",
			PAD+"build"+"\t"+"compile packages and dependencies",
			PAD+"help"+"\t"+"this",
			PAD+"new"+"\t\t"+"create new go module",
			PAD+"run"+"\t\t"+"compile and run Go program",
			PAD+"test"+"\t"+"run Go test packages",
			PAD+"tidy"+"\t"+"add missing and remove unused modules",
			PAD+"version"+"\t"+"print Go version",
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
			if i == len(lines)-1 {
				lineNums = append(lineNums, i)
			}
		}

		for i := range lineNums {
			if i == len(lineNums)-1 {
				break
			}
			pkgHTML = append(pkgHTML, lines[lineNums[i]:lineNums[i+1]])
		}

		for i := len(pkgHTML) - 1; i >= 0; i-- {
			var pkgName string
			var pkgMeta string
			var pkgDesc string
			for ii, pkgLine := range pkgHTML[i] {
				pkgLine = strings.Trim(pkgLine, " \t\n")
				if strings.Contains(pkgLine, "SearchSnippet-header-path") {
					write := false
					for _, ch := range pkgLine {
						if ch == '(' && write == false {
							write = true
							continue
						} else if ch == ')' && write == true {
							write = false
						}
						if write {
							pkgName += string(ch)
						}
					}
				} else if strings.Contains(pkgLine, "published on") {
					write := false
					for _, ch := range pkgLine {
						if ch == '>' && write == false {
							write = true
							continue
						} else if ch == '<' && write == true {
							write = false
						}
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
		in = strings.Trim(in, " \t\n")
		opt, err := strconv.Atoi(in)
		if err != nil || opt >= pkgQueryLimit || opt < 1 {
			return &CmdError{
				Type: CmdAdd,
				Msg: fmt.Sprintf("index '%s' not found\n\nenter an integer value from 1-%d",
					in,
					pkgQueryLimit,
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
	loop:
		for {
			count++
			fmt.Printf(
				"downloading %s %s\r",
				Bold(pkgNames[opt-1]),
				Color("("+strconv.Itoa(count)+"s)", BLUE),
			)
			time.Sleep(1 * time.Second)
			select {
			case out := <-output:
				fmt.Printf("\n%s", out)
				close(output)
				break loop
			default:
			}
		}
	}

	return nil
}

func version() *CmdError {
	versionCmd := exec.Command("go", "version")
	fmt.Println("gopher version 1.1.0")
	fmt.Print(string(Unwrap(versionCmd.Output())))
	return nil
}
func tidy() *CmdError {
	tidyCmd := exec.Command("go", "mod", "tidy")
	output, err := tidyCmd.CombinedOutput()
	if err != nil {
		return &CmdError{
			Type: CmdTidy,
			Msg:  strings.Split(string(output)[:len(output)-1], ": ")[1],
		}
	}
	return nil
}
func test(args ...string) *CmdError {
	args = append([]string{"test", "-v"}, args...)
	runCmd := exec.Command("go", args...)
	output, _ := runCmd.CombinedOutput()

	newOutput := ""
	passes := 0
	totalTests := 0
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "=== RUN  ") {
			totalTests++
			functionName := strings.Split(string(line), " ")[len(strings.Split(string(line), " "))-1]
			newOutput += fmt.Sprintln(
				"==>",
				Color(functionName, BLUE),
			)
		} else if strings.Contains(line, "--- PASS: ") ||
			strings.Contains(line, "--- FAIL: ") {
			lineArr := strings.Split(line, " ")
			time := lineArr[len(lineArr)-1]
			functionName := lineArr[2]
			outcome := lineArr[1]
			before, after, _ := strings.Cut(newOutput, functionName)

			if outcome == "PASS:" {
				passes++
				outcome = Color(" PASS", GREEN)
			} else if outcome == "FAIL:" {
				outcome = Color(" FAIL", RED)
			}

			newOutput = fmt.Sprint(
				before,
				Color(functionName, BLUE),
				time,
				outcome,
				after,
			)
		} else {
			newOutput += fmt.Sprintln(line)
		}
	}
	newOutput = newOutput[:len(newOutput)-2]
	newOutputArr := strings.Split(newOutput, "\n")
	lastLine := newOutputArr[len(newOutputArr)-1]
	lastLineArr := strings.Split(lastLine, "\t")
	finalOutcome := strings.Trim(lastLineArr[0], " ")
	projectName := lastLineArr[1]
	totalTime := lastLineArr[2]

	if finalOutcome == "ok" {
		finalOutcome = Color(fmt.Sprintf(
			"PASS(%d/%d)", passes, totalTests,
		), GREEN)
		newOutput = strings.Join(newOutputArr[:len(newOutputArr)-2], "\n")
	} else if finalOutcome == "FAIL" {
		finalOutcome = Color(fmt.Sprintf(
			"FAIL(%d/%d)", passes, totalTests,
		), RED)
		newOutput = strings.Join(newOutputArr[:len(newOutputArr)-3], "\n")
	}

	fmt.Println(newOutput)
	fmt.Println("\n",
		Bold(projectName)+"("+totalTime+")",
		finalOutcome,
	)

	return nil
}

var PAD string = "    "
var DEFAULT_PREFERENCES = `PkgQueryLimit=10
OpArchPairs=windows,amd64,windows,arm64,linux,amd64,linux,arm64,darwin,amd64,darwin,arm64
PrettyPrint=true
PrettyPrintPreviewLines=3`

type Command int

const (
	CmdNew = iota
	CmdAdd
	CmdHelp
	CmdTidy
	CmdBuild
	CmdConfig
	CmdRun
	CmdTest
	CmdVersion
	CmdInvalid
)

var commandName = map[Command]string{
	CmdNew:     "new",
	CmdAdd:     "add",
	CmdHelp:    "help",
	CmdTidy:    "tidy",
	CmdBuild:   "build",
	CmdConfig:  "config",
	CmdRun:     "run",
	CmdTest:    "test",
	CmdVersion: "version",
	CmdInvalid: "invalid",
}

func (c Command) String() string {
	return commandName[c]
}

type CmdError struct {
	Type Command
	Msg  string
}

func (c CmdError) Error() string {
	return c.Msg
}
