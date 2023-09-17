package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	cmd := ""
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}

	var err error
	switch cmd {
	case "build":
		err = build(os.Args[2:])

	case "serve":
		err = serve(os.Args[2:])

	case "dist":
		err = dist(os.Args[2:])

	case "update":
		err = update(os.Args[2:])

	default:
		usage := `usage: go run ./tool <command> [arguments]

commands:
        build     build in WebAssembly (wasm)
        serve     run a local server
        dist      copy the artifacts to the 'dist' directory
        dist -zip bundle the artifacts as 'dist.zip'
        update    update dependencies and necessary files

tips:
        To modify the contents of the distribution, edit dist.go.
        To modify the build process, edit build.go.
`
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	if err != nil {
		log.Fatalf("%s: %v", cmd, err)
	}
}
