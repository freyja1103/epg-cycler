## EPG-CYCLER

常時起動しないシステムでの EDCB の録画後実行用スクリプト
(ほぼアニメ向け)

### 動作

-   録画したファイルをタイトル別にフォルダを作成し，そのフォルダ内に保存します．
    -   基本半角に変換して保存します．
-   録画後に次の日の朝 4 時までに予約がない場合は，システムをシャットダウンします．
    -   `-process=hoge.exe`のようにプロセス名を指定することで，一致したプロセスが動作中の場合にシャットダウンを防ぎます．

### 使い方

-   EpgTimer 側で HTTP サーバを有効化しておきます．
-   `go build .` or [release](https://github.com/freyja1103/epg-cycler/releases) からダウンロードした後に，

```
chcp 65001
epg-cycler.exe -srcpath="save/path" -originpath=$FilePath$ -title=$TitleF$ -basename=$FileName$ -process="something.exe"
```

のような bat ファイルを作り，EpgTimer 側の録画後実行 bat に設定します．

#### params

| 引数名        | デフォルト値       | 説明                                                                      |
| ------------- | ------------------ | ------------------------------------------------------------------------- |
| `-srcpath`    | `-`                | 保存先の動画パス（動画を保存するディレクトリ）                            |
| `-originpath` | `-`                | 元の動画パス（録画直後などの一時保存パス）                                |
| `-title`      | `-`                | 番組名                                                                    |
| `-basename`   | `-`                | 拡張子なしのファイル名（例：`MyShow_ep1`）                                |
| `-process`    | `-`                | シャットダウンを抑止するプロセス名                                        |
| `-address`    | `"localhost:5510"` | サーバーアドレス（形式：`host:port`、例：`localhost:5510`）               |
| `-all`        | `false`            | ディレクトリ内の録画ファイルを一括整理するモード（true で全整理を有効化） |

#### ほか

-   一括整理用に`-all`オプション作りました．`-savepath`と併せて使ってください．
