## Command `tool`

([English](https://github.com/eihigh/wasmgame/blob/main/tool/README.md) / 日本語)

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

ブラウザゲーム開発に便利な機能を一つにまとめたコマンドです。

どのコマンドも、`go run ./tool build` のようにプロジェクトルートから実行されることを想定しています。

## 使い方

### build
wasm としてビルドします。

serve を実行中にビルドした場合は、サーバーと連携して自動的にブラウザをリロードさせます。（うまく動作しない場合は、手動で一度リロードしてみてください）

### serve
開発用サーバーを立ち上げます。デフォルトでは `http://localhost:8080` でサービスします。URLは `-http` フラグで変更することが可能です。

また、デフォルトでは可能ならば自動でブラウザが立ち上がりますが、不要な場合は `-no-open` フラグを指定して抑制します。

### dist
配布物を `dist` ディレクトリにコピーします。

`-zip` フラグを指定すると、ディレクトリを作成した後、それを `dist.zip` としてアーカイブします。itch.io など投稿サイトにアップロードするのに便利です。

### update
`go.mod` に記載されている依存関係をアップデートします。

## Tips

配布物の内容を修正するには、`tool/dist.go` の `distFiles` を編集してください。

その他、設定や挙動をいじるには、`tool/` 内の .go ファイルを直接編集してください。
