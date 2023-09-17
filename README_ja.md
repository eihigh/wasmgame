# wasmgame

([English](https://github.com/eihigh/wasmgame/blob/main/README.md) / 日本語)

Go と Ebitengine でブラウザゲームを開発し、無料でインターネットに公開するためのテンプレートです。

## 使い方

### テンプレートからリポジトリを作成する
このリポジトリ (github.com/eihigh/wasmgame) は、テンプレートリポジトリです。誰でも内容を複製することで、手軽に開発を始められます。

`Use this template` という緑色のボタンをクリックし、`Create a new repository` を選んでください。

![Create from template](https://github.com/eihigh/wasmgame/assets/44455895/1da9c20e-532c-4585-9953-7f58fb554e38)

リポジトリ作成画面が開くので、複製後の新たなリポジトリ名を入力し、`Create repository` をクリックすると、リポジトリが作成されます。このとき、公開 (Public) か非公開 (Private) を選ぶことができますが、無料アカウントでは非公開リポジトリは GitHub Pages で公開できないので、注意してください。

そのリポジトリを `git clone github.com/<yourname>/<reponame>` で手元にダウンロードすれば、準備は完了です。

### 開発する
`go run ./tool build` を実行することで、プログラムをビルドし、`game.wasm` と `wasm_exec.js` を生成します。しかし、ブラウザゲームは `.exe` ファイルのようにダブルクリックで起動することができません。

`go run ./tool serve` を実行している間はサーバーが起動し、 `http://localhost:8080` にアクセスできるようになるので、このURLをブラウザで開くことでゲームをプレイできます。（`go run ./tool serve` の実行を中断すると、アクセスできなくなります。）`localhost` というURLは、インターネットに公開されない、自分のマシンだけでアクセスできる特殊なURLです。

画像などの素材は `asset` ディレクトリ配下に配置してください。また、ブラウザゲームでは通常のプログラムと異なり、素材を読み込むのに `os.Open` を使うことが出来ないので、`main.go` 内に存在する `open` 関数及び `readFile` 関数を利用して読み込んでください。

サーバーを起動している間に `go run ./tool build` を実行すると、自動的にサーバーと連携して、ブラウザが再読み込みを行います。活用してください。

### GitHub Pages で公開する
作ったゲームは GitHub Pages を利用し、無料でインターネットに公開することができます。ただし、前述の通り無料アカウントでは非公開リポジトリは GitHub Pages を利用できないので、公開リポジトリを使ってください。

まずは、GitHub Pages 機能を有効にします。[⚙ Settings] タブ > サイドバーの [Pages] > [Build and deployment] セクションの、[Source] の選択肢から [GitHub Actions] を選択すると、有効になります。

![GitHub Pages settings](https://github.com/eihigh/wasmgame/assets/44455895/6637c9c0-74f7-4bdc-8c2e-1b2fa950ca98)

有効にしたら、以後 `main` ブランチに git push することで、自動的に GitHub Pages へ公開する処理が始まります。処理が完了したら、`https://<yourname>.github.io/<reponame>` にアクセスし、世界中の誰でも遊べる状態になります。

その他、`tool` の詳しい使い方は、`tool/README.md` を参照してください。

## 注意事項

### 配布物を追加する
デフォルトでは、センシティブなデータを誤って公開するのを防ぐために、`index.html`, `game.html`, `game.wasm`, `wasm_exec.js`, それと `asset` ディレクトリ配下のファイルだけがサーバーを通じて公開されます。

`favicon.ico` など、別のファイルを追加するには、`tool/dist.go` の `distFiles` を編集してください。

### `Layout` とウィンドウサイズ
Ebitengine の仕様上、ブラウザゲームとしてビルドすると、`SetWindowSize` など Ebitengine のウィンドウ制御関数は無効になります。代わりに、HTML上でゲームが画面上のどこに配置されるか制御する必要があります。`index.html` の `<style>` タグ内を参照してください。

基本的に、`main.go` の `Layout` 関数が返す画面の解像度の縦横比とHTML内の表示領域の縦横比が一致するように制御すると、黒帯無しでうまく表示されます。例えば、このリポジトリでは `Layout` が `1280, 720` を返すので、ゲームが表示される領域の縦横比を `aspect-ratio: 16 / 9` に設定しています。
