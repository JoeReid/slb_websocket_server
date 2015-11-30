package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "slb_websocket_server",
	Short: "A light-weight subscribable websocket server",
	Long: `A light-weight, performant websocket server supporting N-many
subscribable message channels based on unique, client provided, subscription keys.

slb_websocket_server was designed to intergrate into the rest of the
super-limit-break live performance tool-chain, but is generic enough to be
useful as a stand-alone websocket messaging system.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.slb_websocket_server.yaml)")
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	// enable ability to specify config file via flag
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("slb_websocket_server.config")
	viper.AddConfigPath("/etc")
	viper.AutomaticEnv() // read in environment variables that match

	initViperDefaults()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// initViperDefaults sets the default config for the application
func initViperDefaults() {
	// ========================================================================
	// server
	// ========================================================================
	viper.SetDefault("server.port", 8080)
	// noGroupAction specifies if a client that subscribes without specifying a
	// group gets all messages or none. This flag can be specified as 'all' or
	// 'none'
	viper.SetDefault("server.noGroupAction", "none")

	// ========================================================================
	// logging
	// ========================================================================
	viper.SetDefault("logging.level", "info")
}
