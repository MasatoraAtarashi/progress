progress - 進捗取得ツール
=======

[![Apache License](http://img.shields.io/badge/license-Apache-blue.svg?style=flat)](LICENSE)

`progress`は進捗を確認することのできるツールです。日報等に使えます。

具体的には、指定したリポジトリに対して自分が行った当日分のコミットを取得することができます。

ユーザ・日付を変更することも可能です。
## Installation
### Homebrew

	brew tap 
	brew install 

### go get
Install

    $ go get 

Update

    $ go get -u 

## Usage
### 進捗を管理したいリポジトリを追加する
$HOME配下に`/.progress.yaml`という名前のファイルを作成してください。
```$HOME/.progress.yaml
unko
```

### 進捗を取得する
    % a
## Options
```
  Options:
  -h,  --help                   print usage and exit
  -n,  --name <yourname>        specify your name
```

## Author
[MasatoraAtarashi](https://github.com/MasatoraAtarashi)
