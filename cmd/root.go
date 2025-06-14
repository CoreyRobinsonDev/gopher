package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "2.2.0"
	rootCmd = &cobra.Command{
		Use: "gopher",
		Short: "Go module manager",
	}
	config = &Config{}
	PAD = "    "
	BLACK = lipgloss.Color("0")
	RED = lipgloss.Color("1")
	GREEN = lipgloss.Color("2")
	YELLOW = lipgloss.Color("3")
	BLUE = lipgloss.Color("4")
	PURPLE = lipgloss.Color("5")
	CYAN = lipgloss.Color("6")
	GRAY = lipgloss.Color("7")
)

func Execute() {
	initConfig()
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(tidyCmd)
}

func initConfig() {
	home := Unwrap(os.UserHomeDir())

	viper.AddConfigPath(home + "/.config/gopher")
	viper.SetConfigType("json")
	viper.SetConfigName("settings")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		var file *os.File
		file, err := os.Create(home + "/.config/gopher/settings.json")
		if err != nil {
			Expect(os.Mkdir(home + "/.config/gopher", 0755))
			file = Unwrap(os.Create(home + "/.config/gopher/settings.json"))
		}
		defer file.Close()
		config = &Config{
			PrettyPrint: true,
			PrettyPrintPreviewLines: 3,
			PkgQueryLimit: 10,
		}
		configBytes := Unwrap(json.MarshalIndent(config, "", "\t"))
		file.Write(configBytes)
	}
	json.Unmarshal(Unwrap(os.ReadFile(home + "/.config/gopher/settings.json")), &config)
}

func Expect(err error) {
	if err != nil {
		fmt.Fprintf(
			os.Stderr, 
			"%s %v\n",
			lipgloss.NewStyle().Foreground(GRAY).Render("gopher:"),
			err,
		)
		os.Exit(1)
	}
}

func Unwrap[T any](result T, err error) T {
	if err != nil {
		fmt.Fprintf(
			os.Stderr, 
			"%s %v\n",
			lipgloss.NewStyle().Foreground(GRAY).Render("gopher:"),
			err,
		)
		os.Exit(1)
	}
	return result
}
