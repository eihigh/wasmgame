// This file is based on github.com/hajimehoshi/wasmserve and modified by eihigh.

// The original license notice:
//
// Copyright 2018 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	tempDir          = filepath.Join(os.TempDir(), "wasmgame")
	errCachingFailed = errors.New("caching failed")
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
	tags := flag.String("tags", "", "Build tags")
	flag.Parse(args)

	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments:", flag.Args())
		flag.Usage()
	}

	if err := execBuild(*tags); err != nil {
		return err
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

func execBuild(tags string) error {
	v, err := goVersion(".")
	if err != nil {
		return err
	}

	f, err := getCachedWasmExecJS(v)
	if errors.Is(err, fs.ErrNotExist) {
		f, err = fetchWasmExecJS(v)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err := os.WriteFile("wasm_exec.js", f, 0666); err != nil {
		return err
	}

	cmd := exec.Command("go", "build", "-o", "game.wasm", "-tags", tags)
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}

	return nil
}

// goVersion fetches the current using Go's version.
// goVersion is different from runtime.Version(), which returns a Go version for this wasmgame build.
func goVersion(builddir string) (string, error) {
	cmd := exec.Command("go", "list", "-f", "go{{.Module.GoVersion}}", builddir)

	var stderr bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s%w", stderr.String(), err)
	}

	return strings.TrimSpace(string(out)), nil
}

func getCachedWasmExecJS(version string) ([]byte, error) {
	cachefile := filepath.Join(tempDir, version, "misc", "wasm", "wasm_exec.js")
	return os.ReadFile(cachefile)
}

func fetchWasmExecJS(version string) ([]byte, error) {
	url := fmt.Sprintf("https://go.googlesource.com/go/+/refs/tags/%s/misc/wasm/wasm_exec.js?format=TEXT", version)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, resp.Body))
	if err != nil {
		return nil, err
	}

	cachefile := filepath.Join(tempDir, version, "misc", "wasm", "wasm_exec.js")
	if err := os.MkdirAll(filepath.Dir(cachefile), 0777); err != nil {
		return content, errors.Join(err, errCachingFailed)
	}
	if err := os.WriteFile(cachefile, content, 0777); err != nil {
		return content, errors.Join(err, errCachingFailed)
	}

	return content, nil
}
