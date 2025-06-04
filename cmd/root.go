package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use: "gopher",
		Short: "Go module manager",
	}
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
		ReportCaller: false,
		Prefix: "gopher",
	})
)

func Execute() {
	Expect(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	home := Unwrap(os.UserHomeDir())

	viper.AddConfigPath(home + "/.config/gopher")
	viper.SetConfigType("json")
	viper.SetConfigName("settings")
	viper.AutomaticEnv()

	Expect(viper.ReadInConfig())
	viper.ConfigFileUsed()
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
