package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ongasatoshi/scion/internal/git"
	"github.com/ongasatoshi/scion/pkg/output"
	"github.com/spf13/cobra"
)

var (
	clearForce      bool
	clearAll        bool
	clearKeepBranch bool
)

var clearCmd = &cobra.Command{
	Use:   "clear <branch-name>",
	Short: "既存のworktreeブランチを削除",
	Long: `clear コマンドは既存のGit Worktreeブランチとその関連ディレクトリを削除します。

例:
  scion clear feature/old-feature
  scion clear feature/experimental --force
  scion clear feature/temp --keep-branch
  scion clear --all`,
	Args: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if all {
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("ブランチ名を指定してください (または --all フラグを使用)")
		}
		return nil
	},
	RunE: runClear,
}

func init() {
	rootCmd.AddCommand(clearCmd)

	clearCmd.Flags().BoolVarP(&clearForce, "force", "f", false, "未コミットの変更があっても強制的に削除")
	clearCmd.Flags().BoolVarP(&clearAll, "all", "a", false, "すべてのworktreeを削除")
	clearCmd.Flags().BoolVar(&clearKeepBranch, "keep-branch", false, "worktreeは削除するがブランチは保持")
}

func runClear(cmd *cobra.Command, args []string) error {
	// Gitリポジトリかどうか確認
	if !git.IsGitRepository() {
		return fmt.Errorf("Gitリポジトリ内で実行してください")
	}

	if clearAll {
		return runClearAll()
	}

	branchName := args[0]
	return clearWorktree(branchName)
}

func runClearAll() error {
	worktrees, err := git.ListWorktrees()
	if err != nil {
		return err
	}

	// メインworktree（最初のエントリ）を除外
	var toRemove []git.WorktreeInfo
	for i, wt := range worktrees {
		if i == 0 || wt.IsBare {
			continue
		}
		toRemove = append(toRemove, wt)
	}

	if len(toRemove) == 0 {
		output.Info("削除するworktreeがありません")
		return nil
	}

	// 削除対象を表示
	output.Info("削除対象のworktree:")
	for _, wt := range toRemove {
		fmt.Printf("  - %s (%s)\n", wt.Branch, wt.Path)
	}

	// 確認プロンプト（--force でない場合）
	config := GetConfig()
	if config.UI.ConfirmDestructive && !clearForce {
		fmt.Print("\nすべてのworktreeを削除しますか? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			output.Info("キャンセルしました")
			return nil
		}
	}

	// 削除を実行
	for _, wt := range toRemove {
		if err := clearWorktreeByPath(wt.Path, wt.Branch); err != nil {
			output.Error("'%s' の削除に失敗しました: %v", wt.Branch, err)
			continue
		}
		output.Success("worktreeを削除しました: %s", wt.Branch)
	}

	output.Success("すべてのworktreeを削除しました")
	return nil
}

func clearWorktree(branchName string) error {
	// リポジトリルートを取得
	repoRoot, err := git.GetRepositoryRoot()
	if err != nil {
		return err
	}

	config := GetConfig()
	baseDir := config.Worktree.BaseDir

	// worktreeパスを構築
	parentDir := filepath.Dir(repoRoot)
	safeBranchName := strings.ReplaceAll(branchName, "/", "-")
	worktreePath := filepath.Join(parentDir, baseDir, safeBranchName)

	return clearWorktreeByPath(worktreePath, branchName)
}

func clearWorktreeByPath(worktreePath, branchName string) error {
	// worktreeが存在するか確認
	if !git.WorktreeExists(worktreePath) {
		return fmt.Errorf("worktree '%s' が見つかりません", worktreePath)
	}

	// 未コミットの変更を確認
	if !clearForce {
		hasChanges, err := git.HasUncommittedChanges(worktreePath)
		if err != nil {
			output.Warning("ステータスの確認に失敗しました: %v", err)
		} else if hasChanges {
			return fmt.Errorf("worktree '%s' には未コミットの変更があります\n--force オプションで強制削除できます", branchName)
		}
	}

	// worktreeを削除
	output.Info("worktreeを削除しています: %s", branchName)
	if err := git.RemoveWorktree(worktreePath, clearForce); err != nil {
		return err
	}
	output.Success("Worktreeを削除しました: %s", worktreePath)

	// ブランチも削除（--keep-branch でない場合）
	if !clearKeepBranch {
		if err := git.DeleteBranch(branchName, clearForce); err != nil {
			output.Warning("ブランチの削除に失敗しました: %v", err)
		} else {
			output.Success("ブランチ '%s' を削除しました", branchName)
		}
	}

	return nil
}
