# config - サブコマンド仕様書

## 概要
`config`コマンドはscionの設定を管理します。設定の表示、更新、リセットなどの機能を提供します。

## 構文
```bash
scion config [subcommand] [flags]
```

## サブコマンド
- `get <key>` - 特定の設定値を取得
- `set <key> <value>` - 設定値を更新
- `list` - すべての設定を表示
- `reset` - 設定をデフォルトに戻す
- `edit` - エディタで設定ファイルを開く

## フラグ
- `--global` - グローバル設定を対象とする
- `--local` - ローカル（リポジトリ固有）設定を対象とする
- `-h, --help` - configコマンドのヘルプを表示

## 設定項目

### グローバル設定（~/.config/scion/config.toml）
```toml
# Worktree関連の設定
[worktree]
base_dir = "wtree"              # worktreeディレクトリのベース名
auto_create_dir = true          # wtreeディレクトリを自動作成
cleanup_on_branch_delete = true # ブランチ削除時にworktreeも削除

# Git関連の設定
[git]
default_remote = "origin"       # デフォルトのリモート名
default_base_branch = "main"    # デフォルトのベースブランチ
fetch_before_create = true      # worktree作成前にfetchを実行

# UI関連の設定
[ui]
color_output = true             # カラー出力を有効化
verbose = false                 # 詳細な出力
confirm_destructive = true      # 破壊的操作の確認

# エディタ設定
[editor]
command = "vi"                  # デフォルトエディタ
```

### ローカル設定（.scion/config.toml）
リポジトリ固有の設定を格納:
```toml
[worktree]
base_dir = "custom-wtree"       # このリポジトリ専用のwtreeディレクトリ名

[git]
default_remote = "upstream"     # このリポジトリのデフォルトリモート
```

## 動作仕様

### 1. config get
```bash
scion config get <key>
```
- ドット記法でネストされた設定にアクセス
- 例: `scion config get worktree.base_dir`
- ローカル設定を優先、なければグローバル設定を参照

### 2. config set
```bash
scion config set <key> <value>
```
- 設定値を更新
- 例: `scion config set worktree.base_dir custom-wtree`
- デフォルトではグローバル設定を更新
- `--local`フラグでローカル設定を更新

### 3. config list
```bash
scion config list
```
- すべての設定を階層的に表示
- ローカル設定がある場合は、それを優先表示
- 設定の出所（global/local）を明示

### 4. config reset
```bash
scion config reset [--all]
```
- 特定の設定またはすべての設定をデフォルトに戻す
- 確認プロンプトを表示（`--force`で省略可能）

### 5. config edit
```bash
scion config edit
```
- 設定ファイルをエディタで開く
- エディタは`$EDITOR`環境変数または設定値を使用
- 保存時に設定の妥当性を検証

## 設定の優先順位
1. コマンドラインフラグ
2. 環境変数（SCION_*）
3. ローカル設定（.scion/config.toml）
4. グローバル設定（~/.config/scion/config.toml）
5. デフォルト値

## 出力例

### get コマンド
```bash
$ scion config get worktree.base_dir
wtree
```

### set コマンド
```bash
$ scion config set worktree.base_dir custom-wtree
✓ Configuration updated: worktree.base_dir = custom-wtree
```

### list コマンド
```bash
$ scion config list
Global Configuration (~/.config/scion/config.toml):
  worktree:
    base_dir: wtree
    auto_create_dir: true
    cleanup_on_branch_delete: true
  git:
    default_remote: origin
    default_base_branch: main
    fetch_before_create: true
  ui:
    color_output: true
    verbose: false
    confirm_destructive: true

Local Configuration (.scion/config.toml):
  worktree:
    base_dir: custom-wtree (overrides global)
```

### reset コマンド
```bash
$ scion config reset
Are you sure you want to reset all configurations to default? [y/N]: y
✓ Configuration reset to defaults
```

## 使用例
```bash
# 特定の設定値を取得
scion config get git.default_remote

# グローバル設定を更新
scion config set ui.color_output false

# ローカル設定を更新
scion config set worktree.base_dir project-wtree --local

# すべての設定を表示
scion config list

# 設定をエディタで編集
scion config edit

# グローバル設定をリセット
scion config reset --global

# 特定の設定項目をリセット
scion config reset worktree.base_dir
```

## エラーハンドリング
- 無効な設定キーの場合はエラーメッセージを表示
- 設定ファイルの構文エラーは詳細な位置情報と共に表示
- 権限不足で設定ファイルが更新できない場合は適切な警告

## 注意事項
- 設定ファイルはTOML形式で管理
- ローカル設定はリポジトリ毎に独立
- 破壊的な変更（reset等）は確認プロンプトを表示
