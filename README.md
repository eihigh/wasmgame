# wasmgame

(English / [日本語](https://github.com/eihigh/wasmgame/blob/main/README_ja.md))

This is a template repository for publishing games developed with Go and Ebitengine as a browser game on the Internet for free.

Your game will run on both Desktop and Browser with just one command!

## Features

* 🚀 Kickstart your development in no time!
* 🛠️ Everything you need, right out of the box!
* 📤 Push to deploy on GitHub Pages – for FREE!

## Demo
https://eihigh.github.io/wasmgame/

(It's very plain and I'd like to improve it...)

## Tutorial
1. Click the green button `Use this template` at top of the page and select `Create a new repository`.
2. Enter a name for your new repository and click `Create repository`.
3. Download the new repository with `git clone`.
4. Go to the directory containing the downloaded `main.go` or `go.mod` and run `go run ./tool build`.
5. Run `go run ./tool serve` and open `http://localhost:8080` in your browser.
6. Verify that the game screen appears.
7. Edit `main.go` as you want and run `go run ./tool build` again in another terminal with `go run ./tool serve` running to see if the browser reloads automatically.
8. Repeat 7.

## Usage

### Creating a repository from a template
This repository (github.com/eihigh/wasmgame) is a template repository. Anyone can easily start development by duplicating this repository.

Click the green button `Use this template` and select `Create a new repository`.

![Create from template](https://github.com/eihigh/wasmgame/assets/44455895/1da9c20e-532c-4585-9953-7f58fb554e38)

Enter the name of the new repository and click "Create repository" to create the repository. At this point, you can choose Public or Private, but be aware that if you have a free account and a private repository, you will not be able to publish pages on GitHub Pages.

Download the new repository to your local machine via `git clone github.com/<yourname>/<reponame>` and you are ready to go.

### How to develop
In addition to the usual Ebitengine development, the following steps will make your game work in the browser!

* Run `go run ./tool build` to build for browsers.
* Run `go run ./tool serve` and open `http://localhost:8080` in your browser to play your game.
* Place assets under `asset` and load them with the `open` function instead of `os.Open`.

Run `go run ./tool build` to build the program and generate `game.wasm` and `wasm_exec.js`. However, browser games cannot be launched with a double-click like `.exe` files.

`go run ./tool serve` will start the server and make `http://localhost:8080` accessible, so you can play the game by opening this URL in your browser. `localhost` is a special URL that is not published on the Internet and can be accessed only on your machine.

Assets such as images should be placed under the `asset` directory. Unlike normal programs, browser games cannot use `os.Open` to read assets, so use the `open` and `readFile` functions in `main.go` to read them. These functions can be used commonly on desktop and browser.

Now your game will also work in the browser! Of course, if you build and run the game normally, it will run on the desktop as well.

While the server is running, `go run ./tool build` will automatically work with the server to reload the browser. Please take advantage of this.

For more information on how to use the `tool`, please refer to `tool/README.md`.

### Publish on GitHub Pages
You can publish your game to the Internet for free using GitHub Pages. However, as mentioned above, you cannot use GitHub Pages for private repositories with a free account, so please use public repositories.

First, enable the GitHub Pages feature. You can do this by going to the [⚙ Settings] tab > [Pages] of the sidebar > [Build and deployment] section and selecting [GitHub Actions] under [Source].

![GitHub Pages settings](https://github.com/eihigh/wasmgame/assets/44455895/6637c9c0-74f7-4bdc-8c2e-1b2fa950ca98)

Once enabled, `git push` to your `main` branch will automatically start the process of publishing to GitHub Pages. When the process is complete, you can go to `https://<yourname>.github.io/<reponame>` and anyone in the world can play with it.

### How to design the page
As you can see in [demo](https://eihigh.github.io/wasmgame/), the game is embedded in a part of the page, and you can freely design outside of it with HTML.

By editing `index.html` in this repository, you can freely design outside of the game. On the other hand, `game.html` is what controls the embedding of the game, so you are unlikely to edit it.

If you know HTML, you can add other files to make it more fancy and link it to your own website. Please also see "Adding Distributions" below.

## Notes

### Adding Distributions
By default, only the `index.html`, `game.html`, `game.wasm`, `wasm_exec.js`, and `asset` directory are published through the server to prevent accidental publication of sensitive data.

To add other files, such as `favicon.ico`, edit `distFiles` in `tool/dist.go`.

### `Layout` and the window size
Due to Ebitengine specifications, Ebitengine's window control functions such as `SetWindowSize` are disabled when built as a browser game. Instead, you need to control where the game is positioned on the screen in HTML, see the `<style>` tag in `index.html`.

Basically, if the aspect ratio of the screen resolution returned by the `Layout` function in `main.go` matches the aspect ratio of the display area in the HTML, the game will display well without black bars. For example, in this repository, `Layout` returns `640, 360`, so the aspect ratio of the element on which the game is displayed is set to `aspect-ratio: 16 / 9`.

### License
This repository is licensed under the BSD Zero Clause License (0BSD). You may omit the original license notice when duplicating or forking this repository.
