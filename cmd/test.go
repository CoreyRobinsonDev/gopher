package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	testCmd = &cobra.Command{
		Use: "test",
		Short: "run _test.go files",
		Long: fmt.Sprintf("run _test.go files\n\n%s %s\n",
			lipgloss.NewStyle().Foreground(YELLOW).Bold(true).Render("example:"),
			lipgloss.NewStyle().Italic(true).Render(
				fmt.Sprintf("gopher %s",
					lipgloss.NewStyle().Foreground(CYAN).Render("test"),
				),
			),
		),
		Run: func(cmd *cobra.Command, args []string) {
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
						lipgloss.NewStyle().Foreground(CYAN).Render(functionName),
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
						outcome = lipgloss.NewStyle().Foreground(GREEN).Render(" PASS")
					} else if outcome == "FAIL:" {
						outcome = lipgloss.NewStyle().Foreground(RED).Render(" FAIL")
					}

					newOutput = fmt.Sprint(
						before,
						lipgloss.NewStyle().Foreground(CYAN).Render(functionName),
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
				finalOutcome = lipgloss.
					NewStyle(). 
					Foreground(GREEN). 
					Render(fmt.Sprintf("PASS(%d/%d)", passes, totalTests))
				newOutput = strings.Join(newOutputArr[:len(newOutputArr)-2], "\n")
			} else if finalOutcome == "FAIL" {
				finalOutcome = lipgloss.
					NewStyle(). 
					Foreground(RED). 
					Render(fmt.Sprintf("FAIL(%d/%d)", passes, totalTests))
				newOutput = strings.Join(newOutputArr[:len(newOutputArr)-2], "\n")
				newOutput = strings.Join(newOutputArr[:len(newOutputArr)-3], "\n")
			}

			fmt.Println(newOutput)
			fmt.Println("\n",
				lipgloss.NewStyle().Bold(true).Render(projectName)+"("+totalTime+")",
				finalOutcome,
				)
		},
	}
)


