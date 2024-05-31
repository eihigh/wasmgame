# wasmgame

([English](https://github.com/eihigh/wasmgame/blob/main/README.md) / 日本語)

Go と Ebitengine で開発したゲームを、ブラウザゲームとして公開するためのテンプレート＋便利ツールです。

## 特長

* 🌏 デスクトップとブラウザ、どちらでも動くゲームが簡単に作れる！
* 🚀 アップするだけで自動でWebページ (GitHub Pages) として公開！
* 📤 itch.io などの投稿サイトに上げる zip も作れる！
* 💲 ぜんぶ無料！

## デモ
こちらが実際に GitHub Pages へ自動で公開されたページです。（とても地味なので、改善したい...）

https://eihigh.github.io/wasmgame/ 

## 始め方
このリポジトリの `Use this template` -> `Create a new repository` をクリック！

![Create from template](https://github.com/eihigh/wasmgame/assets/44455895/1da9c20e-532c-4585-9953-7f58fb554e38)

`main.go` を編集し、通常の Ebitengine プロジェクトと同様にゲームを開発します。ただし、素材は `asset/` ディレクトリ配下に保存し、読み込みには `os.Open` 関数ではなく、main.go にあらかじめ定義された `open` 関数を使ってください。

`go run ./tool build` でブラウザ向けビルドを行い、`go run ./tool serve` で実際にブラウザでプレイして確認することができます。

## 詳しい使い方

### テンプレートからリポジトリを作成する
このリポジトリ (github.com/eihigh/wasmgame) は、テンプレートリポジトリです。複製することで誰でも手軽に、デスクトップとブラウザ両方で動くゲームを開発できます。

`Use this template` という緑色のボタンをクリックし、`Create a new repository` を選んでください。

![Create from template](https://github.com/eihigh/wasmgame/assets/44455895/1da9c20e-532c-4585-9953-7f58fb554e38)

リポジトリ作成画面が開くので、複製後の新たなリポジトリ名を入力します。`Create repository` をクリックすると、リポジトリが作成されます。このとき、公開 (Public) か非公開 (Private) かを選ぶことができますが、**無料アカウントでは非公開リポジトリを GitHub Pages で公開できない**ので、GitHub Pages へ公開したい場合は、Public を選んでください。

そのリポジトリを `git clone github.com/<yourname>/<reponame>` で手元にダウンロードすれば、準備は完了です。

### 開発する
いつもの Ebitengine 開発にプラスして、以下の手順を踏むことで、ブラウザでも動作するゲームになります！

* 画像など素材は `asset/` ディレクトリ配下に保存し、`os.Open` ではなく main.go にあらかじめ定義された `open` 関数で読み込む。
* `go run ./tool build` でブラウザ向けビルドを行う。
* `go run ./tool serve` を実行し、ブラウザで `http://localhost:8080` を開いて動作確認する。
* 以後、serve 中に `go run ./tool build` すると自動でページがリロードされるので便利。

画像などの素材は `asset/` ディレクトリ配下に保存してください。また、ブラウザゲームでは通常のプログラムと異なり、素材を読み込むのに `os.Open` を使うことが出来ないので、main.go にあらかじめ定義された `open` 関数を使ってください。この関数なら、デスクトップでもブラウザでも共通で使えます。

`go build` の代わりに `go run ./tool build` を実行することで、プログラムをビルドし、`game.wasm` と `wasm_exec.js` を生成します。これらのファイルは `.exe` ファイルのようにダブルクリックで起動することができないため、サーバーを起動する必要があります。サーバーを起動するには `go run ./tool serve` を実行します。

`go run ./tool serve` を実行している間はサーバーが起動し、 `http://localhost:8080` にアクセスできるようになるので、このURLをブラウザで開くことでゲームをプレイできます。（`go run ./tool serve` の実行を中断すると、アクセスできなくなります。）`localhost` というURLは、インターネットに公開されない、自分のマシンだけでアクセスできる特殊なURLです。

サーバーを起動している間に `go run ./tool build` を実行すると、自動的にページがリロードされます。ぜひ活用してください。

これで、あなたのゲームがブラウザでも遊べます！ もちろん、普通に `go build` して普通にデスクトップアプリとして遊ぶことも可能です。

その他、`tool` の詳しい使い方は、`tool/README.md` を参照してください。

### GitHub Pages で公開する
作ったゲームは GitHub Pages を利用し、無料でインターネットに公開することができます。ただし、前述の通り無料アカウントでは非公開リポジトリは GitHub Pages を利用できないので、公開リポジトリを使ってください。

まずは、複製したリポジトリの GitHub Pages 機能を有効にします。[⚙ Settings] タブ > サイドバーの [Pages] > [Build and deployment] セクションの、[Source] の選択肢から [GitHub Actions] を選択することで、有効になります。

![GitHub Pages settings](https://github.com/eihigh/wasmgame/assets/44455895/6637c9c0-74f7-4bdc-8c2e-1b2fa950ca98)

有効にしたら、以後 `main` ブランチに git push することで、自動的に GitHub Pages へ公開する処理が始まります。処理が完了したら、`https://<yourname>.github.io/<reponame>` にアクセスし、世界中の誰でも遊べる状態になります。

### ページを装飾する
[デモ](https://eihigh.github.io/wasmgame/)を見ていただくと分かる通り、ページの一部にゲームが埋め込まれ、その外はHTMLで自由にデザインすることが可能です。

このリポジトリの `index.html` を編集することで、ゲームの外を自由に装飾できます。HTMLファイルとしてもう一つ `game.html` がありますが、これはゲームを埋め込むための制御を行っているものなので、編集することはあまりないでしょう。

HTMLの知識があれば、他にもファイルを追加してもっと豪華にしたり、自分のウェブサイトとリンクさせたりできます。後述の「配布物を追加する」も合わせてご確認ください。

## 注意事項

### 配布物を追加する
デフォルトでは、センシティブなデータを誤って公開するのを防ぐために、`index.html`, `game.html`, `game.wasm`, `wasm_exec.js`, それと `asset` ディレクトリ配下のファイルだけがサーバーを通じて公開されます。

`favicon.ico` など、別のファイルを追加するには、`tool/dist.go` の `distFiles` を編集してください。

### `Layout` とウィンドウサイズ
Ebitengine の仕様上、ブラウザゲームとしてビルドすると、`SetWindowSize` など Ebitengine のウィンドウ制御関数は無効になります。代わりに、`index.html` の `<style>` タグなどでゲームが画面上のどこに配置されるか制御する必要があります。

このとき、HTMLで指定したゲーム領域の縦横比と、Ebitengine の Layout メソッドが返す画面解像度の縦横比が一致していないと、黒帯が発生してしまいます。例えば、このリポジトリでは `Layout` が `640, 360` を返すので、ゲームが表示される領域の縦横比を `aspect-ratio: 16 / 9` に設定しています。

### ライセンス
このリポジトリはBSD Zero Clause License (0BSD) の下でライセンスされています。このリポジトリを複製やフォークした際は、元のライセンス表記は省略したり削除してかまいません。

`tool/serve.go` と `tool/build.go` は Apache License 2.0 で頒布されている github.com/hajimehoshi/wasmserve を eihigh が改変したものです。
