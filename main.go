package main

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type game struct {
	sampleJSON []byte
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 64, 64, 255})
	ebitenutil.DebugPrint(screen, "Hello, wasmgame!\nThe content of asset/sample.json is: "+string(g.sampleJSON))
}

func (g *game) Layout(w, h int) (int, int) {
	return 640, 360 // Screen resolution (not window size)
}

func main() {
	g := &game{}
	g.sampleJSON, _ = readFile("asset/sample.json")
	ebiten.SetWindowSize(1280, 720) // has no effect on browser
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

// open opens a file. In a browser, it downloads the file via HTTP;
// otherwise, it reads the file on disk.
func open(name string) (io.ReadCloser, error) {
	name = filepath.Clean(name)
	if runtime.GOOS == "js" {
		// TODO: use more lightweight method such as marwan-at-work/wasm-fetch
		resp, err := http.Get(name)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}

	return os.Open(name)
}

func readFile(name string) ([]byte, error) {
	f, err := open(name)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}
	defer f.Close()

	return io.ReadAll(f)
}
