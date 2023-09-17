package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

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

const defaultAddr = "localhost:8080"

var waitCh = make(chan struct{})

func serve(args []string) error {
	// Parse flags
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
		fmt.Fprintln(os.Stderr, "Unexpected arguments:", flag.Args())
		flag.Usage()
	}

	// Register handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Handle special paths
		switch r.URL.Path {
		case "/_notify":
			notifyForWait(w, r)
			return

		case "/_wait":
			waitForNotify(w, r)
			return
		}

		// Disable caching
		w.Header().Set("Cache-Control", "no-store")

		if *allowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", *allowOrigin)
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
			time.Sleep(time.Second * time.Duration(*delay))
		}

		http.ServeFile(w, r, file)
	})

	log.Println("Listening on", *addr)

	// Open browser if possible.
	if !*noOpen {
		u := "http://" + *addr

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
			log.Printf("Please open %s on your browser.\n", u)
		}
	}

	return http.ListenAndServe(*addr, nil)
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

func waitForNotify(w http.ResponseWriter, r *http.Request) {
	waitCh <- struct{}{}
	http.ServeContent(w, r, "", time.Now(), bytes.NewReader(nil))
}

func notifyForWait(w http.ResponseWriter, r *http.Request) {
	for {
		select {
		case <-waitCh:
		default:
			http.ServeContent(w, r, "", time.Now(), bytes.NewReader(nil))
			return
		}
	}
}
