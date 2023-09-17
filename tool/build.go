package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func build(args []string) error {
	// Parse flags
	flag := flag.NewFlagSet("build", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go run ./tool build [arguments]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	addr := flag.String("http", defaultAddr, "HTTP service address")
	flag.Parse(args)

	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "Unexpected arguments:", flag.Args())
		flag.Usage()
	}

	// Copy $GOROOT/misc/wasm/wasm_exec.js
	goroot := findGOROOT()
	src := filepath.Join(goroot, "misc", "wasm", "wasm_exec.js")
	dst := "wasm_exec.js"
	if err := copyFile(dst, src); err != nil {
		return fmt.Errorf("copy wasm_exec.js: %w", err)
	}

	// Run go build
	cmd := exec.Command("go", "build", "-o", "game.wasm")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}

	// After building, send a request to '_notify' to automatically reload the browser
	u := url.URL{
		Scheme: "http",
		Host:   *addr,
		Path:   "/_notify",
	}

	// Ignore the error, as the build can be done even if the server is not running
	http.PostForm(u.String(), nil)

	return nil
}

func findGOROOT() string {
	if env := os.Getenv("GOROOT"); env != "" {
		return filepath.Clean(env)
	}
	def := filepath.Clean(runtime.GOROOT())
	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return def
	}
	return strings.TrimSpace(string(out))
}
