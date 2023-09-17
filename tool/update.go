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

	log.Println("Update go (go get go)")
	cmd := exec.Command("go", "get", "go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go get go: %w", err)
	}

	log.Println("Update dependencies (go get all)")
	cmd = exec.Command("go", "get", "all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go get all: %w", err)
	}

	log.Println("Cleanup dependencies (go mod tidy)")
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	return nil
}
