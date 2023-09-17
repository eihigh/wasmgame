# wasmgame

(English / 日本語)

This is a template repository for developing browser games with Go and Ebitengine and publishing them on the Internet for free.

## Usage

### Creating a repository from a template
This repository (github.com/eihigh/wasmgame) is a template repository. Anyone can easily start development by duplicating this repository.

Click the green button `Use this template` and select `Create a new repository`.

Enter the name of the new repository and click "Create repository" to create the repository. At this point, you can choose Public or Private, but be aware that if you have a free account and a private repository, you will not be able to publish pages on GitHub Pages.

Download the new repository to your local machine via `git clone github.com/<yourname>/<reponame>` and you are ready to go.

### How to develop
Run `go run . /tool build` to build the program and generate `game.wasm` and `wasm_exec.js`. However, browser games cannot be launched with a double-click like `.exe` files.

`go run . /tool serve` will start the server and make `http://localhost:8080` accessible, so you can play the game by opening this URL in your browser. `localhost` is a special URL that is not published on the Internet and can be accessed only on your machine.

Assets such as images should be placed under the `asset` directory. Unlike normal programs, browser games cannot use `os.Open` to read assets, so use the `open` and `readFile` functions in `main.go` to read them.

While the server is running, `go run . /tool build` will automatically work with the server to reload the browser. Please take advantage of this.

### Publish on GitHub Pages
You can publish your game to the Internet for free using GitHub Pages. However, as mentioned above, you cannot use GitHub Pages for private repositories with a free account, so please use public repositories.

When you `git push` to your `main` branch, the process of publishing to GitHub Pages will start automatically. Once the process is complete, you can access `https://<yourname>.github.io/<reponame>` and anyone in the world can play with it.

For more information on how to use the `tool`, please refer to `tool/README.md`.

## Notes

### Adding Distributions
By default, only the `index.html`, `game.html`, `game.wasm`, `wasm_exec.js`, and `asset` directory is published through the server to prevent accidental publication of sensitive data.

To add other files, such as `favicon.ico`, edit `distFiles` in `tool/dist.go`.

### `Layout` and the window size
Due to Ebitengine specifications, Ebitengine's window control functions such as `SetWindowSize` are disabled when built as a browser game. Instead, you need to control where the game is positioned on the screen in HTML, see the `<style>` tag in `index.html`.

Basically, if the aspect ratio of the screen resolution returned by the `Layout` function in `main.go` matches the aspect ratio of the display area in the HTML, the game will display well without black bars. For example, in this repository, `Layout` returns `1280, 720`, so the aspect ratio of the element on which the game is displayed is set to `aspect-ratio: 16 / 9`.

