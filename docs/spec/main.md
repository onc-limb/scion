# scion - メインコマンド仕様書

## 概要
scionはGit Worktreeを効率的に管理するためのCLIツールです。AI開発環境（Claude Code、GitHub Copilot、Cursor）との連携を前提とした、シンプルで使いやすいワークフロー管理を提供します。

## コマンド構成
```
scion [flags]
scion [command]
```

## 利用可能なコマンド
- `create` - 新しいworktreeブランチを作成
- `clear` - 既存のworktreeブランチを削除
- `config` - scionの設定を管理

## グローバルフラグ
- `-h, --help` - ヘルプ情報を表示
- `-v, --version` - バージョン情報を表示
- `--config string` - 設定ファイルのパスを指定（デフォルト: `~/.config/scion/config.toml`）

## 設定ファイル
### 場所
- デフォルトパス: `~/.config/scion/config.toml`
- `go install`実行時に自動的に作成される

### 設定ファイル形式
TOML形式で以下の設定を管理:
```toml
# scion設定ファイル

[repository]
base_repository = "" # コマンド対象のリポジトリパス
base_branch = "main" # worktreeのブランチを切るベースブランチ名

[worktree]
base_dir = "wtree"  # worktreeディレクトリのベース名

[git]
default_remote = "origin"  # デフォルトのリモート名
```

## 初期化処理
1. `go install`実行時に設定ディレクトリを確認
2. 設定ファイルが存在しない場合は、デフォルト設定で作成
3. 必要な権限の確認と設定

## エラーハンドリング
- Gitリポジトリ外での実行時はエラーメッセージを表示
- 設定ファイルの読み込み失敗時は、デフォルト値を使用
- 各サブコマンドのエラーは適切にユーザーに通知

## 使用例
```bash
# ヘルプの表示
scion --help

# バージョンの確認
scion --version

# 新しいworktreeの作成
scion create feature/new-feature

# worktreeの削除
scion clear feature/old-feature

# 設定の更新
scion config set worktree.base_dir custom-wtree
```
