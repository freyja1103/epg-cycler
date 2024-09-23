## EPG-CYCLER

EDCB の録画後実行用スクリプト

### 動作

-   録画したファイルをタイトル別にフォルダを作成し，そのフォルダ内に保存します．
-   録画後に次の日の朝 4 時までに録画がない場合は，PC をシャットダウンします．
    -   `-process=hoge.exe`のようにプロセス名を指定することでシャットダウンを防ぎます．

### 使い方
- EpgTimer側でHTTPサーバを有効化しておきます．
- `go build`した後に，

```
epg-cycler.exe -title=$SCtitle$ -subtitle=$SCsubtitle$ -number=$SCnumber$ -process=something.exe
```

のような bat ファイルを作り，EpgTimer 側の録画後実行 bat に設定します．
