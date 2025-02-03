package main

import (
    "bytes"
    "os"
    "testing"

    "github.com/spf13/viper"
)

func TestExecuteNoArgs(t *testing.T) {
    buf := new(bytes.Buffer)
    rootCmd.SetOut(buf)
    rootCmd.SetErr(buf)

    os.Args = []string{"tfbackend"}

    if err := Execute(); err != nil {
        t.Fatalf("rootCmd.Execute() failed: %v", err)
    }

    output := buf.String()
    if len(output) == 0 {
        t.Fatalf("Expected output, got none")
    }
}

func TestConfigInitialization(t *testing.T) {
    cfgFile = ""
    os.Clearenv()
    os.Setenv("HOME", "/tmp")

    initConfig()

    if viper.GetString("listen_addr") != ":8080" {
        t.Fatalf("Expected listen_addr to be :8080, got %s", viper.GetString("listen_addr"))
    }
}
