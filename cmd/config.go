package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use: "config",
		Short: "configure your gopher settings",
		Long: "configure your gopher settings\n  this file can be found at ~/.config/gopher/settings.json",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

type Config struct {
	PrettyPrint bool `json:"prettyPrint"`
	PrettyPrintPreviewLines uint `json:"prettyPrintPreviewLines"`
	PkgQueryLimit uint `json:"pkgQueryLimit"`
}
