package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	cfg, err := LoadConfig("../")
	if err != nil {
		slog.Error("init failed:", "error", err)
		os.Exit(0)
	}

	fmt.Println(cfg)
}
