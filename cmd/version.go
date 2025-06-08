package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use: "version",
		Short: "print go and gopher version",
		Run: func(cmd *cobra.Command, args []string) {
			versionCmd := exec.Command("go", "version")
			fmt.Println("gopher version v" + version)
			fmt.Print(string(Unwrap(versionCmd.Output())))
		},
	}
)

