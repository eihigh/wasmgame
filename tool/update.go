package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func update(args []string) error {
	// update can only be runnable in the project root
	if _, err := os.Stat("go.mod"); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("workdir is not the project root")
	}

	log.Println("Reset dependencies")
	if err := os.Remove("go.sum"); err != nil {
		return fmt.Errorf("Remove go.sum: %w", err)
	}
	b := []byte("module example.com/game")
	if err := os.WriteFile("go.mod", b, 0666); err != nil {
		return fmt.Errorf("Rewrite go.mod: %w", err)
	}

	log.Println("Update dependencies (go mod tidy)")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	return nil
}
