package cmd

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	PAD = "    "
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
	config = &Config{
		PrettyPrint: true,
		PrettyPrintPreviewLines: 3,
		PkgQueryLimit: 10,
	}
)

func Execute() {
	Expect(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(addCmd)
}

func initConfig() {
	home := Unwrap(os.UserHomeDir())

	viper.AddConfigPath(home + "/.config/gopher")
	viper.SetConfigType("json")
	viper.SetConfigName("settings")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		file := Unwrap(os.Create(home + "/.config/gopher/settings.json"))
		defer file.Close()
		configBytes := Unwrap(json.MarshalIndent(config, "", "\t"))
		file.Write(configBytes)
	}
	json.Unmarshal([]byte(viper.ConfigFileUsed()), &config)
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
