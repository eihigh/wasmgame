package main

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	fetch "marwan.io/wasm-fetch"
)

type game struct {
	ticks      int
	sampleJSON []byte
}

func (g *game) Update() error {
	g.ticks++
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 64, 64, 255})
	s := fmt.Sprintf("Hello, wasmgame!\nTicks = %d\nThe content of asset/sample.json is: %s", g.ticks, string(g.sampleJSON))
	x, y := g.ticks%640, g.ticks%360
	ebitenutil.DebugPrintAt(screen, s, x, y)
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
		resp, err := fetch.Fetch(name, &fetch.Opts{})
		if err != nil {
			return nil, err
		}
		return io.NopCloser(bytes.NewReader(resp.Body)), nil
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
