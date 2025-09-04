/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/namnd/xai-cli/cmd/chat"
	"github.com/namnd/xai-cli/local"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xai",
	Short: "xAI wrapper CLI for coding productivity",
	Long: `A WIP CLI tool, built on top of xAI models (grok-3-mini, grok-code-fast-1 & grok-4)

It provides basic features to understand context of a codebase.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(initConfig)

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xai-cli.yaml)")

	rootCmd.AddCommand(chat.ChatCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, local.CONFIG_DIR)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		// config.yaml file not found, create new one with defaults
		if err := os.MkdirAll(configDir, 0700); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create config directory: %v\n", err)
			os.Exit(1)
		}

		viper.SetDefault("model", "grok-3-mini")
		viper.SetDefault("timeout", 120)

		configFile := filepath.Join(configDir, "config.yaml")
		if err := viper.WriteConfigAs(configFile); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create intial config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Initial config created ", configDir)
	}

	viper.SetEnvPrefix("XAI_CLI")
	viper.AutomaticEnv()
}
