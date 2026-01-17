# clear - サブコマンド仕様書

## 概要
`clear`コマンドは既存のGit Worktreeブランチとその関連ディレクトリを削除します。

## 構文
```bash
scion clear [flags] <branch-name>
```

## 引数
- `<branch-name>` (必須) - 削除するworktreeブランチ名

## フラグ
- `-f, --force` - 未コミットの変更があっても強制的に削除
- `-a, --all` - すべてのworktreeを削除
- `--keep-branch` - worktreeは削除するがブランチは保持
- `-h, --help` - clearコマンドのヘルプを表示

## 動作仕様

### 1. 前提条件の確認
- Gitリポジトリ内で実行されているか確認
- 指定されたworktreeが存在するか確認
- 削除対象が現在作業中のworktreeでないか確認

### 2. 処理フロー
1. worktreeの存在確認
   ```bash
   git worktree list
   ```
2. 指定されたブランチ名のworktreeを特定
3. 未コミットの変更を確認
   - 変更あり + `--force`フラグなし: 警告を表示して処理を中断
   - 変更あり + `--force`フラグあり: 処理を続行
   - 変更なし: 処理を続行
4. worktreeを削除
   ```bash
   git worktree remove <worktree-path>
   ```
5. `wtree`ディレクトリから対象ディレクトリを削除
6. `--keep-branch`フラグがない場合、ブランチも削除
   ```bash
   git branch -d <branch-name>
   ```
7. 削除完了メッセージを表示

### 3. 特殊な動作

#### 全worktree削除（--allフラグ）
```bash
scion clear --all
```
1. すべてのworktreeをリストアップ（メインworktreeを除く）
2. 確認プロンプトを表示
3. ユーザーの確認後、すべてのworktreeを順次削除

#### ブランチ保持（--keep-branchフラグ）
```bash
scion clear feature/temp --keep-branch
```
- worktreeディレクトリは削除
- Gitブランチは保持（後で再利用可能）

### 4. エラーケース
- 指定されたworktreeが存在しない場合
- 現在作業中のworktreeを削除しようとした場合
- 未コミットの変更があり、`--force`フラグがない場合
- 削除権限がない場合
- Git Worktreeコマンドが失敗した場合

### 5. 出力例

#### 成功時
```bash
$ scion clear feature/old-feature
Removing worktree: feature/old-feature
✓ Worktree removed: ../wtree/feature/old-feature
✓ Branch 'feature/old-feature' deleted
```

#### 警告時（未コミットの変更あり）
```bash
$ scion clear feature/working
Warning: Worktree 'feature/working' has uncommitted changes
Use --force to remove anyway, or commit/stash your changes first
```

#### 全削除時
```bash
$ scion clear --all
Found 3 worktrees to remove:
  - feature/login
  - feature/payment
  - bugfix/issue-123

Are you sure you want to remove all worktrees? [y/N]: y
✓ Removed worktree: feature/login
✓ Removed worktree: feature/payment
✓ Removed worktree: bugfix/issue-123
All worktrees cleared successfully
```

## 使用例
```bash
# 基本的な使用
scion clear feature/old-feature

# 未コミットの変更があっても強制削除
scion clear feature/experimental --force

# worktreeは削除するがブランチは保持
scion clear feature/temp --keep-branch

# すべてのworktreeを削除
scion clear --all

# すべてのworktreeを強制削除（確認なし）
scion clear --all --force
```

## 安全性の考慮事項
- デフォルトでは未コミットの変更がある場合は削除を拒否
- 現在作業中のworktreeは削除できない
- `--all`フラグ使用時は確認プロンプトを表示
- 削除前にworktreeの状態をログに記録

## 注意事項
- 削除されたworktreeは復元できない
- `--keep-branch`を使用しない限り、ブランチも一緒に削除される
- メインworktree（元のリポジトリ）は削除対象外
