# create - サブコマンド仕様書

## 概要
`create`コマンドは新しいGit Worktreeブランチを作成し、専用のディレクトリに配置します。

## 構文
```bash
scion create [flags] <branch-name>
```

## 引数
- `<branch-name>` (必須) - 作成するブランチ名

## フラグ
- `-b, --base string` - ベースブランチを指定（デフォルト: 現在のブランチ）
- `-r, --remote string` - リモートリポジトリを指定（デフォルト: origin）
- `-f, --force` - 既存のブランチを強制的に上書き
- `-h, --help` - createコマンドのヘルプを表示

## 動作仕様

### 1. 前提条件の確認
- Gitリポジトリ内で実行されているか確認
- Git Worktreeが利用可能か確認
- 指定されたブランチ名の妥当性を検証

### 2. ディレクトリ構造
```
<repository-root>/
├── .git/
├── src/
└── ../wtree/           # ベースリポジトリと同一階層
    └── <branch-name>/  # worktreeディレクトリ
```

### 3. 処理フロー
1. 現在のリポジトリルートを取得
2. ベースリポジトリの親ディレクトリに移動
3. `wtree`ディレクトリの存在確認
   - 存在しない場合: `wtree`ディレクトリを作成
   - 存在する場合: そのまま続行
4. Git Worktreeコマンドを実行
   ```bash
   git worktree add ../wtree/<branch-name> -b <branch-name>
   ```
5. 作成成功メッセージを表示
6. 新しいworktreeのパスを出力

### 4. エラーケース
- ブランチ名が既に存在する場合
  - `--force`フラグなし: エラーメッセージを表示して終了
  - `--force`フラグあり: 既存のworktreeを削除して再作成
- `wtree`ディレクトリの作成に失敗した場合
- Git Worktreeコマンドが失敗した場合
- 権限不足でディレクトリ作成ができない場合

### 5. 出力例

#### 成功時
```bash
$ scion create feature/new-feature
Creating worktree for branch: feature/new-feature
✓ Worktree directory created: ../wtree/feature/new-feature
✓ Branch 'feature/new-feature' created and checked out
Path: /path/to/repo/../wtree/feature/new-feature
```

#### エラー時
```bash
$ scion create feature/existing-feature
Error: Branch 'feature/existing-feature' already exists
Use --force to overwrite existing worktree
```

## 使用例
```bash
# 基本的な使用
scion create feature/login

# ベースブランチを指定
scion create feature/payment --base develop

# リモートブランチから作成
scion create bugfix/issue-123 --remote upstream

# 既存のworktreeを強制的に再作成
scion create feature/refactor --force
```

## 注意事項
- worktreeディレクトリ（`wtree`）は、ベースリポジトリと同一階層に作成される
- ブランチ名にはGitの命名規則が適用される
- worktree内での作業は、元のリポジトリに影響を与える（同じGitデータベースを共有）
- 設定の優先度はフラグによる指定>configファイルの指定>デフォルト値とする。
