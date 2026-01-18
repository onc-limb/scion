package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsGitRepository は現在のディレクトリがGitリポジトリ内かどうかを確認する
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// GetRepositoryRoot はGitリポジトリのルートパスを返す
func GetRepositoryRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Gitリポジトリのルートを取得できません: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch は現在のブランチ名を返す
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("現在のブランチを取得できません: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// BranchExists はブランチが存在するかどうかを確認する
func BranchExists(branchName string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branchName)
	return cmd.Run() == nil
}

// WorktreeExists はworktreeが存在するかどうかを確認する
func WorktreeExists(path string) bool {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "worktree ") {
			wtPath := strings.TrimPrefix(line, "worktree ")
			if wtPath == absPath {
				return true
			}
		}
	}
	return false
}

// CreateWorktree は新しいworktreeを作成する
func CreateWorktree(path, branchName, baseBranch string, force bool) error {
	// ディレクトリを作成
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("worktreeディレクトリの作成に失敗しました: %w", err)
	}

	// git worktree add コマンドを構築
	args := []string{"worktree", "add"}

	if force {
		args = append(args, "--force")
	}

	args = append(args, path)

	if BranchExists(branchName) {
		// 既存ブランチをチェックアウト
		args = append(args, branchName)
	} else {
		// 新規ブランチを作成
		args = append(args, "-b", branchName)
		if baseBranch != "" {
			args = append(args, baseBranch)
		}
	}

	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("worktreeの作成に失敗しました: %s", stderr.String())
	}

	return nil
}

// RemoveWorktree はworktreeを削除する
func RemoveWorktree(path string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, path)

	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("worktreeの削除に失敗しました: %s", stderr.String())
	}

	return nil
}

// DeleteBranch はブランチを削除する
func DeleteBranch(branchName string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}

	cmd := exec.Command("git", "branch", flag, branchName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ブランチの削除に失敗しました: %s", stderr.String())
	}

	return nil
}

// ListWorktrees はすべてのworktreeをリストアップする
func ListWorktrees() ([]WorktreeInfo, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("worktreeのリスト取得に失敗しました: %w", err)
	}

	var worktrees []WorktreeInfo
	var current *WorktreeInfo

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "worktree ") {
			if current != nil {
				worktrees = append(worktrees, *current)
			}
			current = &WorktreeInfo{
				Path: strings.TrimPrefix(line, "worktree "),
			}
		} else if strings.HasPrefix(line, "branch ") && current != nil {
			branch := strings.TrimPrefix(line, "branch refs/heads/")
			current.Branch = branch
		} else if line == "bare" && current != nil {
			current.IsBare = true
		}
	}

	if current != nil {
		worktrees = append(worktrees, *current)
	}

	return worktrees, nil
}

// WorktreeInfo はworktreeの情報を保持する
type WorktreeInfo struct {
	Path   string
	Branch string
	IsBare bool
}

// HasUncommittedChanges は未コミットの変更があるかどうかを確認する
func HasUncommittedChanges(worktreePath string) (bool, error) {
	cmd := exec.Command("git", "-C", worktreePath, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("ステータスの確認に失敗しました: %w", err)
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// Fetch はリモートから最新の情報を取得する
func Fetch(remote string) error {
	cmd := exec.Command("git", "fetch", remote)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("fetchに失敗しました: %s", stderr.String())
	}

	return nil
}
