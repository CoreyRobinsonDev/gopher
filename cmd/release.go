package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	CleanFlag bool
	releaseCmd = &cobra.Command{
		Use: "release",
		Short: "release your software for remote deployments",
		Run: func(cmd *cobra.Command, args []string) {
			releaseExe, err := exec.LookPath("goreleaser")
			if err != nil {
				var shouldInstall bool
				confirm := huh.NewConfirm(). 
					Title(fmt.Sprintf("\033[0mModule release utility %s (%s) not found. Would you like to install it?",
						lipgloss.NewStyle().Foreground(CYAN).Render("goreleaser"),
						lipgloss.NewStyle().Foreground(BLUE).Render("https://github.com/goreleaser/goreleaser"),
						)). 
					Affirmative("yes"). 
					Negative("no"). 
					Value(&shouldInstall). 
					WithTheme(huh.ThemeCatppuccin())
				Expect(confirm.Run())

				if shouldInstall {
					goCmd := exec.Command("go", "install", "github.com/goreleaser/goreleaser/v2@latest")
					fmt.Print(string(Unwrap(goCmd.CombinedOutput())))
					fmt.Printf("%s installed successfully in $GOBIN\n",
						lipgloss.NewStyle().Foreground(CYAN).Render("goreleaser"),
					)
				}

				os.Exit(0)
			}
			// NOTE: What is needed to run:
			// - a commit has to be made
			// - a remote repo is set
			// - a git tag is set for the commit hash
			// - local code must be pushed to remote
			if CleanFlag {
				Expect(syscall.Exec(releaseExe, []string{"goreleaser", "release", "--clean"}, os.Environ()))
			} else {
				Expect(syscall.Exec(releaseExe, []string{"goreleaser", "release"}, os.Environ()))
			}
		},
	}
)

func init() {
	releaseCmd.PersistentFlags().BoolVarP(&CleanFlag, "clean", "c", false, "removes the 'dist' directory")	
}
