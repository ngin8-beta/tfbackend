package main

import (
    "os"
    "log/slog"
)

func main() {
    if err := Execute(); err != nil {
        slog.Error("Command execution failed", "error", err)
        os.Exit(1)
    }
}
