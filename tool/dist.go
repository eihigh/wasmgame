package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const distRoot = "dist"

var distFiles = []string{
	"index.html",
	"game.html",
	"game.wasm",
	"wasm_exec.js",
}

var distDirs = []string{
	"asset",
}

// isDist reports whether the name is part of the distribution.
func isDist(name string) bool {
	name = filepath.Clean(name)

	// Return true if the name exactly matches distFiles.
	for _, f := range distFiles {
		f = filepath.Clean(f)
		if name == f {
			return true
		}
	}

	// Skip files that begin with a dot and do not match distFiles.
	base := filepath.Base(name)
	if strings.HasPrefix(base, ".") {
		return false
	}

	// Return true if the file is under distDirs.
	for _, d := range distDirs {
		d = filepath.Clean(d)
		if strings.HasPrefix(name, d) {
			return true
		}
	}

	return false
}

func dist(args []string) error {
	// Parse flags
	flag := flag.NewFlagSet("dist", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go run ./tool dist [arguments]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	zips := flag.Bool("zip", false, "bundle the artifacts as dist.zip")
	flag.Parse(args)

	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "Unexpected arguments:", flag.Args())
		flag.Usage()
	}

	// Build before copying
	if err := build(nil); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	// Reset the existing distribution
	if err := os.RemoveAll(distRoot); err != nil {
		return fmt.Errorf("remove existing dist: %w", err)
	}
	if err := os.MkdirAll(distRoot, 0777); err != nil {
		return fmt.Errorf("ensure dist: %w", err)
	}

	// Copy files
	for _, f := range distFiles {
		dst := filepath.Join(distRoot, f)
		d := filepath.Dir(dst)
		if err := os.MkdirAll(d, 0777); err != nil {
			return fmt.Errorf("ensure dir for %s: %w", f, err)
		}
		if err := copyFile(dst, f); err != nil {
			return fmt.Errorf("copy file %s: %w", f, err)
		}
	}

	// Copy directories recursively
	for _, d := range distDirs {
		dst := filepath.Join(distRoot, d)
		if err := copyDir(dst, d); err != nil {
			return fmt.Errorf("copy dir %s: %w", d, err)
		}
	}

	// Zip dist if needed
	if *zips {
		if err := zipDist(); err != nil {
			return fmt.Errorf("zip dist: %w", err)
		}
	}

	return nil
}

func copyDir(dst, src string) error {
	return filepath.WalkDir(src, func(name string, entry fs.DirEntry, err error) error {
		// Abort the entire copyDir if the entry has an error
		if err != nil {
			return fmt.Errorf("walk %s: %w", name, err)
		}

		// Check if the entry is part of the distribution
		if !isDist(name) {
			log.Println("skipped:", name)
			return nil
		}

		// Ensure directory if the entry is a directory
		dst := filepath.Join(distRoot, name)
		if entry.IsDir() {
			if err := os.MkdirAll(dst, 0777); err != nil {
				return fmt.Errorf("ensure dir for %s: %w", name, err)
			}
			return nil
		}

		// Copy the file
		if err := copyFile(dst, name); err != nil {
			return fmt.Errorf("copy file %s: %w", name, err)
		}
		return nil
	})
}

func copyFile(dst, src string) error {
	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}

func zipDist() error {
	output, err := os.Create("dist.zip")
	if err != nil {
		return fmt.Errorf("create dist.zip: %w", err)
	}
	defer output.Close()

	zw := zip.NewWriter(output)
	defer zw.Close()

	// Write contents of the distribution to 'dist.zip' recursively
	return filepath.WalkDir(distRoot, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk %s: %w", name, err)
		}
		if entry.IsDir() {
			return nil
		}

		f, err := os.Open(name)
		if err != nil {
			return fmt.Errorf("open %s: %w", name, err)
		}
		defer f.Close()

		rel, _ := filepath.Rel(distRoot, name)
		w, err := zw.Create(rel)
		if err != nil {
			return fmt.Errorf("create %s as %s in zip: %w", name, rel, err)
		}

		if _, err = io.Copy(w, f); err != nil {
			return fmt.Errorf("copy %s as %s into zip: %w", name, rel, err)
		}
		return nil
	})
}
