package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use: "config",
		Short: "configure your gopher settings",
		Long: fmt.Sprintf("configure your gopher settings. This file can be found at ~/.config/gopher/settings.json\n\n%s %s",
			lipgloss.NewStyle().Foreground(YELLOW).Bold(true).Render("example:"),
			lipgloss.NewStyle().Italic(true).Render(
				fmt.Sprintf("gopher %s",
					lipgloss.NewStyle().Foreground(CYAN).Render("config"),
				),
			),
		),	
		Run: func(cmd *cobra.Command, args []string) {
			prettyPrintPreviewLines := fmt.Sprintf("%d", config.PrettyPrintPreviewLines)
			pkgQueryLimit := fmt.Sprintf("%d", config.PkgQueryLimit)
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[bool]().
						TitleFunc(func() string {
							return fmt.Sprintf(
								"\033[0mEnable pretty printing on errors when running %s and %s\n\033[0m\033[2mprettyPrint: %v\033[0m",
								lipgloss.NewStyle().Foreground(CYAN).Bold(true).Render("gopher run"),
								lipgloss.NewStyle().Foreground(CYAN).Bold(true).Render("gopher build"),
								config.PrettyPrint,
								)
					}, &config.PrettyPrint).
						Options(
							huh.NewOption("true", true),
							huh.NewOption("false", false),
						).Value(&config.PrettyPrint),
					huh.NewInput().
						TitleFunc(func() string {
							return fmt.Sprintf(
								"\033[0mHow many lines to show before and after an error when prettyPrint is enabled\n\033[0m\033[2mprettyPrintPreviewLines: %v\033[0m",
								prettyPrintPreviewLines,
								)
					}, &prettyPrintPreviewLines).
						Value(&prettyPrintPreviewLines).
						Validate(func (in string) error {
							_, err := strconv.ParseUint(in, 10, 0)
							return err
						}),
					huh.NewInput().
						TitleFunc(func() string {
							return fmt.Sprintf(
								"\033[0mHow many packages to fetch when calling %s\n\033[0m\033[2mpkgQueryLimit: %v\033[0m",
								lipgloss.NewStyle().Foreground(CYAN).Bold(true).Render("gopher add"),
								pkgQueryLimit,
								)					
						}, &pkgQueryLimit).
						Value(&pkgQueryLimit).
						Validate(func (in string) error {
							val, err := strconv.ParseUint(in, 10, 0)
							if val > 100 {
								return errors.New("package query limit cannot be over 100")
							}
							return err
						}),
				),		
			).WithTheme(huh.ThemeCatppuccin())
			Expect(form.Run())
			pppl := Unwrap(strconv.ParseUint(prettyPrintPreviewLines, 10, 0))
			pql := Unwrap(strconv.ParseUint(pkgQueryLimit, 10, 0))
			config.PrettyPrintPreviewLines = uint(pppl)
			config.PkgQueryLimit = uint(pql)
			home := Unwrap(os.UserHomeDir())
			Expect(os.WriteFile(home +"/.config/gopher/settings.json", Unwrap(json.MarshalIndent(config, "", "\t")), 0655))
		},
	}
)

type Config struct {
	PrettyPrint bool `json:"prettyPrint"`
	PrettyPrintPreviewLines uint `json:"prettyPrintPreviewLines"`
	PkgQueryLimit uint `json:"pkgQueryLimit"`
}
