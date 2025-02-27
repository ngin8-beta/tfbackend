package main

import (
	"log/slog"
	"os"

	internalServer "github.com/ngin8-beta/tfbackend/internal/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "tfbackend",
    Short: "Terraform state backend server",
    Long:  "tfbackend is an HTTP server that provides Terraform state management and locking mechanisms. It offers flexible configuration through environment variables and configuration files, making it highly operable.",
	Run: runServer,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to the configuration file (e.g., $HOME/.tfbackend.yaml)")
	rootCmd.PersistentFlags().StringP("listen_port", "p", "8080", "The port to listen on.")
	rootCmd.PersistentFlags().StringP("storage", "s", "local", "The storage backend to use .")

	homedir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to retrieve the home directory.", "error", err)
		os.Exit(1)
	}
	rootCmd.PersistentFlags().StringP("storage_local_dir", "d", homedir + "/tfstate", "The path to the local storage directory. (Required for the local storage backend)")


	if err := viper.BindPFlag("listen_port", rootCmd.PersistentFlags().Lookup("listen_port")); err != nil {
		slog.Error("listen_port Failed to bind the flag.", "error", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("storage", rootCmd.PersistentFlags().Lookup("storage")); err != nil {
		slog.Error("storage Failed to bind the flag.", "error", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("storage_local_dir", rootCmd.PersistentFlags().Lookup("storage_local_dir")); err != nil {
		slog.Error("storage_local_dir Failed to bind the flag.", "error", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile == "" {
		// Use .tfbackend.yaml in the home directory as the default configuration file.
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("Failed to retrieve the home directory.", "error", err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tfbackend")
	} else {
        // Use the configuration file specified on the command line.
		viper.SetConfigFile(cfgFile)
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

func runServer(cmd *cobra.Command, args []string) {
	listenPort := viper.GetString("listen_port")
	storageType := viper.GetString("storage")

	slog.Info("Starting the tfbackend server...", "listen_port", listenPort, "storage", storageType)

	srv := internalServer.NewServer(listenPort)

	if err := srv.Run(); err != nil {
		slog.Error("Failed to start the tfbackend server.", "error", err)
		os.Exit(1)
	}
}