## EPG-CYCLER

常時起動しないシステムでの EDCB の録画後実行用スクリプト
※動作テスト中

### 動作

-   録画したファイルをタイトル別にフォルダを作成し，そのフォルダ内に保存します．
    -   基本半角に変換して保存します．
-   録画後に次の日の朝 4 時までに予約がない場合は，システムをシャットダウンします．
    -   `-process=hoge.exe`のようにプロセス名を指定することで，一致したプロセスが動作中の場合にシャットダウンを防ぎます．

### 使い方

-   EpgTimer 側で HTTP サーバを有効化しておきます．
-   `go build`した後に，

```
chcp 65001
epg-cycler.exe -srcpath="save/path" -originpath=$FilePath$ -title=$TitleF$ -basename=$FileName$ -process="something.exe"
```

のような bat ファイルを作り，EpgTimer 側の録画後実行 bat に設定します．

-   IP とポート番号を指定する場合は，実行時に`-ip=192.168.0.2:7777`のように追記してください．
    -   デフォルトの値は`localhost:5510`です．
-   `-process`は必須ではありません．
-   `chcp 65001`は文字化け回避のために記述してください．
