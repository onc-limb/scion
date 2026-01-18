package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIsGitRepository(t *testing.T) {
	// 現在のディレクトリがGitリポジトリかどうかを確認
	// このテストはGitリポジトリ内で実行されることを前提としている
	result := IsGitRepository()
	if !result {
		t.Skip("Not running inside a git repository")
	}
}

func TestIsGitRepositoryInNonGitDir(t *testing.T) {
	// 一時ディレクトリを作成
	tmpDir := t.TempDir()

	// 一時ディレクトリに移動
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Gitリポジトリでないことを確認
	result := IsGitRepository()
	if result {
		t.Error("expected IsGitRepository to return false in non-git directory")
	}
}

func TestGetRepositoryRoot(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not running inside a git repository")
	}

	root, err := GetRepositoryRoot()
	if err != nil {
		t.Fatalf("failed to get repository root: %v", err)
	}

	// ルートディレクトリに.gitが存在することを確認
	gitDir := filepath.Join(root, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		t.Fatalf(".git not found in repository root: %v", err)
	}

	// .gitはディレクトリまたはファイル（worktreeの場合）である
	if !info.IsDir() {
		// worktreeの場合は.gitはファイル
		data, err := os.ReadFile(gitDir)
		if err != nil {
			t.Fatalf("failed to read .git file: %v", err)
		}
		if len(data) == 0 {
			t.Error(".git file is empty")
		}
	}
}

func TestGetCurrentBranch(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not running inside a git repository")
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		t.Fatalf("failed to get current branch: %v", err)
	}

	if branch == "" {
		t.Error("expected non-empty branch name")
	}
}

func TestBranchExists(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not running inside a git repository")
	}

	// 現在のブランチは存在するはず
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		t.Fatalf("failed to get current branch: %v", err)
	}

	if !BranchExists(currentBranch) {
		t.Errorf("expected current branch '%s' to exist", currentBranch)
	}

	// 存在しないブランチ
	if BranchExists("non-existent-branch-12345") {
		t.Error("expected non-existent branch to return false")
	}
}

func setupTestGitRepo(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Git リポジトリを初期化
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// ユーザー設定
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	cmd.Run()

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	cmd.Run()

	// 初期コミット
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git add: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git commit: %v", err)
	}

	return tmpDir
}

func TestListWorktrees(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not running inside a git repository")
	}

	worktrees, err := ListWorktrees()
	if err != nil {
		t.Fatalf("failed to list worktrees: %v", err)
	}

	// 少なくとも1つのworktree（メインリポジトリ）が存在するはず
	if len(worktrees) == 0 {
		t.Error("expected at least one worktree")
	}
}

func TestHasUncommittedChanges(t *testing.T) {
	tmpDir := setupTestGitRepo(t)

	// クリーンな状態では変更なし
	hasChanges, err := HasUncommittedChanges(tmpDir)
	if err != nil {
		t.Fatalf("failed to check uncommitted changes: %v", err)
	}
	if hasChanges {
		t.Error("expected no uncommitted changes in clean repo")
	}

	// ファイルを変更
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	// 変更があることを確認
	hasChanges, err = HasUncommittedChanges(tmpDir)
	if err != nil {
		t.Fatalf("failed to check uncommitted changes: %v", err)
	}
	if !hasChanges {
		t.Error("expected uncommitted changes after modifying file")
	}
}
