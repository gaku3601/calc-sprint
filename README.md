# 導入方法
https://github.com/gaku3601/calc-spinrt/releases  
から各環境に適した実行ファイルをダウンロードする  
## mac
calc-sprintをダウンロードして.zshrcなどで定義しているpathの通っているところへ配置する(/usr/local/bin等)  
以下コマンドで実行権限を与える  
```
chmod 777 calcspinrt
```

calcspinrt -hでヘルプが表示されれば利用可能です  
セキュリティ周りでpopupが出る場合はmacのセキュリティープライバシーの設定からcalcspinrtを許可してください　　

## win
calcspinrt.exeをダウンロードしてpathの通っているところへ配置する  

# 使い方
```
calcspinrt -p [対象エクセルファイル]
```
を実行することでjiraで生成したExcelファイルからSprintの集計で必要なデータを算出します

## debug

```
go run main.go -p ./example/ex1.xlsx
```

## build

```
cd script
sh build.sh

```
でscript/build.shを回せば、distフォルダにバイナリファイルが格納される
