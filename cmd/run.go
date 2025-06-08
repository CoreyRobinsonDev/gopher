package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use: "run",
		Short: "compile and run Go program",
		Long: fmt.Sprintf("compile and run Go program\n\n%s should be executed at the root of your module\n\n%s %s\n",
			lipgloss.NewStyle().Bold(true).Render("gopher run"),
			lipgloss.NewStyle().Foreground(YELLOW).Bold(true).Render("example:"),
			lipgloss.NewStyle().Italic(true).Render(
				fmt.Sprintf("gopher %s",
					lipgloss.NewStyle().Foreground(CYAN).Render("run"),
				),
			),
		),
		Run: func(cmd *cobra.Command, args []string) {
			buildCmd := exec.Command("go", "build", "-o", "gobinary")
			output, e := buildCmd.CombinedOutput()
			if e == nil {
				Expect(os.Remove("./gobinary"))
			}

			if config.PrettyPrint && e != nil {
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
					previewLines := int(config.PrettyPrintPreviewLines)
					defer f.Close()
					reader := bufio.NewScanner(f)

					fmt.Printf("[%s]\n", lipgloss.NewStyle().Foreground(CYAN).Render(file))
					linenum := 0
					for reader.Scan() {
						linenum++
						if dif := linenum - rownum; dif <= previewLines && dif >= -previewLines {
							tabs := strings.Count(reader.Text(), "\t")
							if linenum == rownum {
								fmt.Printf("%s %s%s\n",
									lipgloss.NewStyle().Foreground(GRAY).Render(strconv.Itoa(linenum)),
									strings.Repeat("  ", tabs),
									lipgloss.NewStyle().Italic(true).Render(strings.Trim(reader.Text(), " \t")),
								)
							} else {
								fmt.Printf("%s %s%s\n",
									lipgloss.NewStyle().Foreground(GRAY).Render(strconv.Itoa(linenum)),
									strings.Repeat("  ", tabs),
									strings.Trim(reader.Text(), " \t"),
									)
							}
						}
						if linenum == rownum {
							tabs := strings.Count(reader.Text(), "\t")
							pad := strings.Repeat(" ", colnum+len(strconv.Itoa(linenum))-tabs)
							fmt.Printf("%s%s%s\n", 
								strings.Repeat("  ", tabs), 
								pad, 
								lipgloss.NewStyle().Foreground(RED).Render("^ "+err),
							)
						}
						if linenum > rownum+previewLines {
							fmt.Println()
							break
						}
					}
				}
			} else if e != nil {
				fmt.Print(string(output))
			} else {
				o := Unwrap(exec.LookPath("go"))
				Expect(syscall.Exec(o, append([]string{"go","run","."}, args...), os.Environ()))
			}
		},
	}
)


