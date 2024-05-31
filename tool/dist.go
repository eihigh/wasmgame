package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const outputDir = "dist"

// distFiles and distDirs are the files and directories to be included in the distribution.
// If you want to include more files or directories, add them to these lists.
var distFiles = []string{
	"index.html",
	"game.html",
	"game.wasm",
	"wasm_exec.js",
	// "favicon.ico",
	// ...
}

// Note that files that start with a dot are excluded even if they are under distDir.
// To include such files, add them to distFiles individually.
var distDirs = []string{
	"asset",
}

// isDist reports whether the name is part of the distribution.
func isDist(name string) bool {
	name, _ = filepath.Abs(name)

	// Return true if the name exactly matches distFiles.
	for _, f := range distFiles {
		f, _ = filepath.Abs(f)
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
		d, _ = filepath.Abs(d)
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
	tags := flag.String("tags", "", "Build tags")
	flag.Parse(args)

	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments:", flag.Args())
		flag.Usage()
	}

	// Build before copying
	if err := execBuild(*tags); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	// Reset the existing output
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("remove existing dist: %w", err)
	}
	if err := os.MkdirAll(outputDir, 0777); err != nil {
		return fmt.Errorf("ensure dist: %w", err)
	}

	// Copy distFiles
	for _, f := range distFiles {
		dst := filepath.Join(outputDir, f)
		if err := copyFile(dst, f); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("files listed in distFiles are required: %s", f)
			} else {
				return fmt.Errorf("copy file %s: %w", f, err)
			}
		}
	}

	// Copy directories recursively
	for _, d := range distDirs {
		dst := filepath.Join(outputDir, d)
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

		rel, err := filepath.Rel(src, name)
		if err != nil {
			return err
		}
		dst := filepath.Join(dst, rel)

		// Ensure directory if the entry is a directory
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
	// We ignore the permissions
	return nil
}

func zipDist() error {
	output, err := os.Create("dist.zip")
	if err != nil {
		return fmt.Errorf("create dist.zip: %w", err)
	}
	defer output.Close()

	zw := zip.NewWriter(output)
	zw.AddFS(os.DirFS(outputDir))
	return zw.Close()
}
