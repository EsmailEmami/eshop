package cmd

import (
	"fmt"
	"os"

	"github.com/esmailemami/eshop/consts"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

var RootCmd = &cobra.Command{
	Use: "billing",
}

func init() {
	cobra.OnInitialize(InitConfig)
	RootCmd.SetHelpTemplate(consts.AppLogo + "\r\n\r\n" + "Version: " + consts.AppVersion + "\r\n\r\n" + RootCmd.HelpTemplate())
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "make output more verbose")
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.yaml)")

	RootCmd.AddCommand(MigrateCmd)
	RootCmd.AddCommand(DBCmd)
}

func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println("")
		os.Exit(1)
	}
}

func InitConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/billing")
		viper.AddConfigPath("$HOME/.billing")
	}

	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
