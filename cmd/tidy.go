package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	tidyCmd = &cobra.Command{
		Use: "tidy",
		Short: "add missing and remove unused modules",
		Long: fmt.Sprintf("add missing and remove unused modules\n\n%s %s\n",
			lipgloss.NewStyle().Foreground(YELLOW).Bold(true).Render("example:"),
			lipgloss.NewStyle().Italic(true).Render(
			fmt.Sprintf("gopher %s",
				lipgloss.NewStyle().Foreground(CYAN).Render("tidy"),
				),
			),
		),
		Run: func(cmd *cobra.Command, args []string) {
			tidyCmd := exec.Command("go", "mod", "tidy")
			output, err := tidyCmd.CombinedOutput()
			if err != nil {
				fmt.Fprintf(
					os.Stderr, 
					"%s %s\n",
					lipgloss.NewStyle().Foreground(GRAY).Render("gopher:"),
					strings.Split(string(output)[:len(output)-1], ": ")[1],
				)
				os.Exit(1)
			}
		},
	}
)

