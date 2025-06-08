package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "2.0.0"
	rootCmd = &cobra.Command{
		Use: "gopher",
		Short: "Go module manager",
	}
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
		ReportCaller: false,
		Prefix: "gopher",
	})
	config = &Config{}
	PAD = "    "
	CYAN = lipgloss.Color("36")
	YELLOW = lipgloss.Color("#ffa500")
	GREEN = lipgloss.Color("#00f000")
	RED = lipgloss.Color("#ff0000")
	GRAY = lipgloss.Color("2")
)

func Execute() {
	initConfig()
	Expect(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(runCmd)
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
		logger.Fatal(err)
	}
}

func Unwrap[T any](result T, err error) T {
	if err != nil {
		logger.Fatal(err)
	}
	return result
}
