package main

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfigNoCfgFile(t *testing.T) {
    // Arrange
    cfgFile = ""

    // Act
    initConfig()

    // Assert
    if viper.GetString("listen_addr") != ":8080" {
        t.Fatalf("Expected listen_addr to be :8080, got %s", viper.GetString("listen_addr"))
    }
}

func TestInitConfigSetCfgFile(t *testing.T) {
    // Arrange
    cfgFile = "/tmp/.tfbackend.yml"
    cfgFileContent := "listen_addr: :3000\n"
    if err := os.WriteFile(cfgFile, []byte(cfgFileContent), 0644); err != nil {
        t.Fatalf("Failed to write file: %v", err)
    }
    defer os.Remove(cfgFile)

    // Act
    initConfig()

    // Assert
    if viper.GetString("listen_addr") != ":3000" {
        t.Fatalf("Expected listen_addr to be :3000, got %s", viper.GetString("listen_addr"))
    }
}

func TestInitConfigSetEnvListenAddr(t *testing.T) {
    // Arrange
    os.Setenv("TFBACKEND_LISTEN_ADDR", ":4000")

    // Act
    initConfig()

    // Assert
    if viper.GetString("listen_addr") != ":4000" {
        t.Fatalf("Expected listen_addr to be :4000, got %s", viper.GetString("listen_addr"))
    }
}

func TestInitConfigSetCfgFileAndEnvListenAddr(t *testing.T) {
    // Arrange
    cfgFile = "/tmp/.tfbackend.yml"
    cfgFileContent := "listen_addr: :3000\n"
    if err := os.WriteFile(cfgFile, []byte(cfgFileContent), 0644); err != nil {
        t.Fatalf("Failed to write file: %v", err)
    }
    defer os.Remove(cfgFile)

    os.Setenv("TFBACKEND_LISTEN_ADDR", ":4000")

    // Act
    initConfig()

    // Assert
    if viper.GetString("listen_addr") != ":4000" {
        t.Fatalf("Expected listen_addr to be :4000, got %s", viper.GetString("listen_addr"))
    }
}