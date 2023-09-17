## Command `tool`

(English / [日本語](https://github.com/eihigh/wasmgame/blob/main/tool/README_ja.md))

```
$ go run ./tool
usage: go run ./tool <command> [arguments]

commands:
        build     build in WebAssembly (wasm)
        serve     run a local server
        dist      copy the artifacts to the 'dist' directory
        dist -zip bundle the artifacts as 'dist.zip'
        update    update dependencies and necessary files

tips:
        To modify the contents of the distribution, edit dist.go.
        To modify the build process, edit build.go.
```

Command `tool` contains all the useful features for browser game development.

All commands are intended to be run from the project root.

## Usage

### build
Build in WebAssembly (wasm).

If you build while the server is running, it will work with the server to automatically reload the browser. (If this does not work, try reloading once manually.)

### serve
Starts the development server. By default, it serves at `http://localhost:8080`. The URL can be changed with the `-http` flag.

The default is to automatically launch the browser if possible, but this can be suppressed with the `-no-open` flag if it is not needed.

### dist
Copies the distribution to the `dist` directory.

The `-zip` flag creates a directory and archives it as `dist.zip`, useful for uploading to sites like itch.io.

### update
Updates the dependencies listed in `go.mod`.

## Tips

To modify the contents of the distribution, edit the `distFiles` in `tool/dist.go`.

To modify other settings or behavior, edit the .go files in `tool/` directly.
