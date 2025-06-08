package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use: "add {package | github.com/user/package}",
		Short: "add dependencies to current module and install them",
		Long: fmt.Sprintf(
			"Add dependencies to current module and install them.\n\nWhen a full package name isn't provided %s will do a search on pkg.go.dev for matching packages.\nThe number of results returned on this search can be adjusted with %s\n\n%s %s\n\t %s",
			lipgloss.NewStyle().Foreground(CYAN).Render("gopher add"),
			lipgloss.NewStyle().Foreground(CYAN).Render("gopher config"),
			lipgloss.NewStyle().Foreground(YELLOW).Bold(true).Render("example:"),
			lipgloss.NewStyle().Italic(true).Render(
				fmt.Sprintf("gopher %s rsc.io/quote",
					lipgloss.NewStyle().Foreground(CYAN).Render("add"),
				),
			),
			lipgloss.NewStyle().Italic(true).Render(
				fmt.Sprintf("gopher %s gofiber",
					lipgloss.NewStyle().Foreground(CYAN).Render("add"),
				),
			),
		),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintf(
					os.Stderr, 
					"%s %s\n",
					lipgloss.NewStyle().Foreground(GRAY).Render("gopher:"),
					"missing package\nrun 'gopher add -h' for usage",
				)
				os.Exit(1)
			}
			pkg := args[0]
			if strings.Contains(pkg, "/") {
				getCmd := exec.Command("go", "get", pkg)
				out, _ := getCmd.CombinedOutput()
				fmt.Print(string(out))
			} else {
				horizonalCharLimit := 80
				url := fmt.Sprintf(
					"https://pkg.go.dev/search?limit=%d&m=package&q=%s",
					config.PkgQueryLimit,
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
						lipgloss.NewStyle().Foreground(CYAN).Render(strconv.Itoa(i+1)),
						lipgloss.NewStyle().Bold(true).Render(pkgName),
						lipgloss.NewStyle().Foreground(CYAN).Render(version),
						lipgloss.NewStyle().Bold(true).Render("("+strings.Join(pkgMetaArr[1:], " ")+")"),
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
				fmt.Printf("%s Packages to install (eg: 1 2 3)\n",
					lipgloss.NewStyle().Foreground(YELLOW).Render("==> "),
				)
				fmt.Print(lipgloss.NewStyle().Foreground(YELLOW).Render("==> "))
				reader := bufio.NewReader(os.Stdin)
				in := Unwrap(reader.ReadString('\n'))
				in = strings.Trim(in, " \t\n")
				opt, err := strconv.Atoi(in)
				if err != nil || opt >= int(config.PkgQueryLimit) || opt < 1 {
					fmt.Fprintf(
						os.Stderr, 
						"%s %s\n",
						lipgloss.NewStyle().Foreground(GRAY).Render("gopher:"),
						fmt.Sprintf("index '%s' not found\n\nenter an integer value from 1-%d",
							in,
							config.PkgQueryLimit,
						),
					)
					os.Exit(1)
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
						lipgloss.NewStyle().Bold(true).Render(pkgNames[opt-1]),
						lipgloss.NewStyle().Foreground(CYAN).Render("("+strconv.Itoa(count)+"s)"),
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
		},
	}
)
