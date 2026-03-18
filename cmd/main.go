package main

import (
	"fmt"
	"os"

	"github.com/somaz94/git-mirror-action/internal/config"
	"github.com/somaz94/git-mirror-action/internal/mirror"
	"github.com/somaz94/git-mirror-action/internal/output"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "::error::%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	m := mirror.New(cfg)
	results := m.Run()

	return output.Write(results)
}
