# wasmgame

(English / [日本語](https://github.com/eihigh/wasmgame/blob/main/README_ja.md))

wasmgame is a template and toolkit that makes it easy to develop games using Go and Ebitengine and publish them as browser games.

## Features

* 🌏 Easily create games that run on both desktop and browser!
* 🚀 Automatically publish your game as a webpage (GitHub Pages) just by uploading it!
* 📤 Generate a zip file to submit your game to platforms like itch.io!
* 💲 Completely free!

## Demo
Here's an example of a page that was automatically published to GitHub Pages using wasmgame. (It's quite basic at the moment, so I'd like to improve it...)

https://eihigh.github.io/wasmgame/

## Getting Started 
Click `Use this template` -> `Create a new repository` in this repository to get started!

![Create from template](https://github.com/eihigh/wasmgame/assets/44455895/1da9c20e-532c-4585-9953-7f58fb554e38)

Edit `main.go` and develop your game just like a regular Ebitengine project. However, make sure you place your assets in the `asset/` directory and use the predefined `open` function in `main.go` to load them instead of `os.Open`.

Run `go run ./tool build` to build for the browser, and `go run ./tool serve` to test it out in your browser.

## Detailed Usage

### Creating a repository from the template
This repository (github.com/eihigh/wasmgame) is a template. By duplicating it, anyone can easily develop games that run on both desktop and browser.

Click the green `Use this template` button and select `Create a new repository`.

![Create from template](https://github.com/eihigh/wasmgame/assets/44455895/1da9c20e-532c-4585-9953-7f58fb554e38) 

On the Create Repository screen, enter a name for your new repository after duplication. Click `Create repository` to create the repo. You can choose to make it public or private, but **free accounts cannot publish private repositories to GitHub Pages**. So if you want to publish to GitHub Pages, choose Public.

Download that repository to your machine with `git clone github.com/<yourname>/<reponame>` and you're all set!

### Developing Your Game  
In addition to the usual Ebitengine development flow, follow these steps to make your game run in the browser too:

* Place assets like images in the `asset/` directory and load them using the predefined `open` function in `main.go` instead of `os.Open`.
* Run `go run ./tool build` to build for the browser. 
* Run `go run ./tool serve` and open `http://localhost:8080` in your browser to test it.
* While serving, running `go run ./tool build` will automatically reload the page which is very handy.

Make sure to put assets like images under the `asset/` directory. Also, unlike regular programs, browser games cannot use `os.Open` to load assets. So use the predefined `open` function in `main.go` which works for both desktop and browser.

Instead of `go build`, run `go run ./tool build` to build your program and generate `game.wasm` and `wasm_exec.js`. These files cannot be launched by double-clicking like an `.exe` file, so you need to start a server. Run `go run ./tool serve` to start the server.

While `go run ./tool serve` is running, you can access `http://localhost:8080`. Open this URL in your browser to play the game. Note that if you stop running `go run ./tool serve`, it will no longer be accessible. The `localhost` URL is a special URL that is not publicly accessible on the internet and can only be accessed from your own machine. 

If you run `go run ./tool build` while the server is running, the page will automatically reload.

Now your game is ready to be played in the browser! Of course, you can also just `go build` and play it as a regular desktop app too.

For more details on how to use the `tool`, see `tool/README.md`.

### Publishing to GitHub Pages
You can use GitHub Pages to publish your game to the web for free. However, as mentioned above, free accounts cannot use GitHub Pages for private repositories, so use a public repository.

First, enable the GitHub Pages feature for your duplicated repository. Go to the [⚙ Settings] tab > [Pages] in the sidebar > Under the [Build and deployment] section, select [GitHub Actions] from the [Source] dropdown to enable it.

![GitHub Pages settings](https://github.com/eihigh/wasmgame/assets/44455895/6637c9c0-74f7-4bdc-8c2e-1b2fa950ca98)

Once enabled, pushing to the `main` branch will automatically trigger the process of publishing to GitHub Pages. When the process is complete, access `https://<yourname>.github.io/<reponame>` and anyone in the world can play your game.

### Customizing the Page
As you can see from the [demo](https://eihigh.github.io/wasmgame/), the game is embedded in part of the page, and the rest can be freely designed with HTML.  

Edit the `index.html` in this repository to freely customize the area outside the game. There's another HTML file called `game.html`, but this handles embedding the game, so you probably won't need to edit it.

If you know HTML, you can add more files to make it fancier and link it to your own website. See "Adding Distribution Files" below for more information.

## Notes

### Adding Distribution Files 
By default, only `index.html`, `game.html`, `game.wasm`, `wasm_exec.js`, and files in the `asset` directory are published through the server to avoid accidentally publishing sensitive data.

To add other files such as `favicon.ico`, edit the `distFiles` in `tool/dist.go`.

### `Layout` and Window size
Due to limitations in Ebitengine, window control functions such as `SetWindowSize` are disabled when running as a browser game. Instead, you have to control where the game is positioned on the screen, for example using the `<style>` tag in `index.html`.

In this case, if the aspect ratio of the game area specified in HTML does not match the aspect ratio of the screen resolution returned by Ebitengine's Layout method, black bars will appear. For example, this repository's `Layout` returns `640, 360`, so the aspect ratio of the area where the game is displayed is set to `aspect-ratio: 16 / 9`.

### License
This repository is licensed under the BSD Zero Clause License (0BSD). When duplicating or forking this repository, you may omit or delete the original license notice.

`tool/serve.go` and `tool/build.go` are modified version of the code from github.com/hajimehoshi/wasmserve licensed under the Apache License 2.0.
