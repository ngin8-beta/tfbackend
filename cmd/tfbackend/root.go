package main

import (
    "os"
    "log/slog"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "tfbackend",
    Short: "Terraform state backend server",
    Long:  "tfbackend is an HTTP server that provides Terraform state management and locking mechanisms. It offers flexible configuration through environment variables and configuration files, making it highly operable.",
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to the configuration file (e.g., $HOME/.tfbackend.yaml)")
	rootCmd.PersistentFlags().String("listen_addr", ":8080", "The address the server listens on.")

	if err := viper.BindPFlag("listen_addr", rootCmd.PersistentFlags().Lookup("listen_addr")); err != nil {
		slog.Error("listen_addr Failed to bind the flag.", "error", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
        // Use the configuration file specified on the command line.
		viper.SetConfigFile(cfgFile)
	} else {
		// Use .tfbackend.yaml in the home directory as the default configuration file.
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("Failed to retrieve the home directory.", "error", err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tfbackend")
	}

    // Environmet variables are prefixed with TFBACKEND_.
    viper.SetEnvPrefix("TFBACKEND")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Configuration file loaded.", "file", viper.ConfigFileUsed())
	} else {
		slog.Info("Configuration file not found. Using default settings and environment variables.")
	}
}
