progress - 進捗取得ツール
=======

[![Apache License](http://img.shields.io/badge/license-Apache-blue.svg?style=flat)](LICENSE)

`progress`は進捗を確認することのできるツールです。日報等に使えます。

具体的には、指定したリポジトリに対して自分が行った当日分のコミットを取得することができます。

ユーザ・日付を変更することも可能です。

## Status
`progress`は指定した日時のslack上での発言を収集する機能とともに[nippo-generater](https://github.com/MasatoraAtarashi/nippo-generator)に統合されます。
## Installation

    $ go get -u github.com/MasatoraAtarashi/progress

## Usage
### 進捗を管理したいリポジトリを追加する
$HOME配下に`/.progress.yaml`という名前のファイルを作成してください(--configオプションで任意の設定ファイルを指定することも可能)。
```$HOME/.progress.yaml
repositories:
    #absolute path to repository you want to manage progress
    - "Users/MasatoraAtarashi/workspace/hogehoge"
```

### 進捗を取得する
    % progress get [options]

#### 実行結果例
```
# 2021-04-25

## progress(3 commits)
 - af24524 Modify description of flag
 - fabf1c8 Add time option
 - fc69ade Add branch option

## dotfiles(2 commits)
 - ea891cb Change file structure
 - fcfce6c Add goland setting
```
## Options
```
  Options:
  -b, --branch <branch-name>   Specify branch
  -d, --date <date>            Specify date like <2021-04-24>
  -h, --help                   help for get
  -r, --reverse                Reverse order of commits
  -t, --time                   Show time
  -u, --user <username>        Specify user
```

## Author
[MasatoraAtarashi](https://github.com/MasatoraAtarashi)
