package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	WebFlag bool
	WatchFlag bool
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
			// build binary first to own the error formatting
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
				if WebFlag {
					goExe := Unwrap(exec.LookPath("go"))
					Expect(syscall.Exec(goExe, append([]string{"go","run", "github.com/hajimehoshi/wasmserve@latest","."}, args...), os.Environ()))
				} else if WatchFlag {
					airExe, err := exec.LookPath("air")
					if err != nil {
						var shouldInstall bool
						confirm := huh.NewConfirm(). 
							Title(fmt.Sprintf("\033[0mLive-reloading utility %s (%s) not found. Would you like to install it? %s",
								lipgloss.NewStyle().Foreground(CYAN).Render("air"),
								lipgloss.NewStyle().Foreground(BLUE).Render("https://github.com/air-verse/air"),
								lipgloss.NewStyle().Foreground(GRAY).Render("(go 1.23 or higher required)"),
							)). 
							Affirmative("yes"). 
							Negative("no"). 
							Value(&shouldInstall). 
							WithTheme(huh.ThemeCatppuccin())
						Expect(confirm.Run())

						if shouldInstall {
							goCmd := exec.Command("go", "install", "github.com/air-verse/air@latest")
							fmt.Print(string(Unwrap(goCmd.CombinedOutput())))
							fmt.Printf("%s installed successfully in $GOBIN\n",
								lipgloss.NewStyle().Foreground(CYAN).Render("air"),
							)
						}

						os.Exit(0)
					}
					Expect(syscall.Exec(airExe, append([]string{"air"}, args...), os.Environ()))
				} else {
					goExe := Unwrap(exec.LookPath("go"))
					Expect(syscall.Exec(goExe, append([]string{"go","run","."}, args...), os.Environ()))
				}
			}
		},
	}
)

func init() {
	runCmd.PersistentFlags().BoolVar(&WebFlag, "web", false, "run program in browser")	
	runCmd.PersistentFlags().BoolVarP(&WatchFlag, "watch", "w", false, "live-reload your code on change")	
}
