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
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const defaultAddr = "localhost:8080"

const reloadScript = `
<script>
(async () => {
  // The server sends a response for '_wait' when a request is sent to '_notify'.
  const reload = await fetch('_wait');
  if (reload.ok) {
    location.reload();
  }
})();
</script>
`

var (
	waitChannel = make(chan struct{})
)

type server struct {
	http.Server
	delay       int
	allowOrigin string
}

func openBrowser(addr string) {
	u := "http://" + addr

	ok := func() bool {
		var err error
		switch runtime.GOOS {
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
		case "darwin":
			err = exec.Command("open", u).Start()
		case "linux":
			err = exec.Command("xdg-open", u).Start()
		default:
			return false
		}

		if err != nil {
			return false
		}
		return true
	}()

	if !ok {
		log.Printf("Please open %s on your browser.", u)
	}
}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	if s.allowOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", s.allowOrigin)
	}

	// Handle special paths
	switch r.URL.Path {
	case "/_notify":
		notifyWaiters(w, r)
		return

	case "/_wait":
		waitForUpdate(w, r)
		return
	}

	// Disable caching
	w.Header().Set("Cache-Control", "no-store")

	// Redirect "/dir" to "/dir/"
	if !strings.HasSuffix(r.URL.Path, "/") {
		fi, err := os.Stat(filepath.Join(".", r.URL.Path))
		if err != nil && !os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if fi != nil && fi.IsDir() {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusSeeOther)
			return
		}
	}

	// Serve files
	file, err := convertPath(r.URL.Path)
	if err != nil {
		log.Printf("%s\t->\t[error]", r.URL.Path)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("%s\t->\t%s", r.URL.Path, file)
	f, err := os.Open(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer f.Close()

	// Inject reload system into index.html
	if filepath.Base(file) == "index.html" {
		b, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Append the reload script after the original contents
		b = append(b, []byte(reloadScript)...)
		http.ServeContent(w, r, "index.html", time.Now(), bytes.NewReader(b))
		return
	}

	// Delay when checking the display of the loading UI
	if strings.HasSuffix(file, ".wasm") {
		time.Sleep(time.Second * time.Duration(s.delay))
	}

	http.ServeFile(w, r, file)
}

// convertPath converts a path of a URL into a file path on the disk.
func convertPath(path string) (string, error) {
	path = filepath.Clean(path)
	path = filepath.Join(".", path)

	stat, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if stat.IsDir() {
		path = filepath.Join(path, "index.html")
	}

	if !isDist(path) {
		return "", fmt.Errorf("%s is not part of the distribution", path)
	}

	return path, nil
}

func waitForUpdate(w http.ResponseWriter, r *http.Request) {
	waitChannel <- struct{}{}
	http.ServeContent(w, r, "", time.Now(), bytes.NewReader(nil))
}

func notifyWaiters(w http.ResponseWriter, r *http.Request) {
	for {
		select {
		case <-waitChannel:
		default:
			http.ServeContent(w, r, "", time.Now(), bytes.NewReader(nil))
			return
		}
	}
}

func serve(args []string) error {
	flag := flag.NewFlagSet("serve", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go run ./tool serve [arguments]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	delay := flag.Int("delay", 0, "Delay for displaying a loading UI")
	addr := flag.String("http", defaultAddr, "HTTP service address")
	allowOrigin := flag.String("allow-origin", "*", "Allowed origin for CORS requests")
	noOpen := flag.Bool("no-open", false, "Do not open browser automatically")
	flag.Parse(args)

	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments:", flag.Args())
		flag.Usage()
	}

	var server server

	shutdown := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Printf("Shutting down server...")

		// Received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Error at server.Shutdown: %v", err)
		}
		close(shutdown)

		<-sigint
		// Hard exit on the second ctrl-c.
		os.Exit(0)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handle)
	server.Handler = mux
	server.Addr = *addr
	server.delay = *delay
	server.allowOrigin = *allowOrigin

	if !(*noOpen) {
		openBrowser(*addr)
	}

	log.Printf("Listening on http://%v", *addr)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-shutdown

	log.Printf("Exiting")
	return nil
}
