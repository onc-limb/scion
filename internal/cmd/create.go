package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ongasatoshi/scion/internal/git"
	"github.com/ongasatoshi/scion/pkg/output"
	"github.com/spf13/cobra"
)

var (
	createBaseBranch string
	createRemote     string
	createForce      bool
)

var createCmd = &cobra.Command{
	Use:   "create <branch-name>",
	Short: "新しいworktreeブランチを作成",
	Long: `create コマンドは新しいGit Worktreeブランチを作成し、専用のディレクトリに配置します。

worktreeディレクトリはベースリポジトリと同一階層に作成されます。

例:
  scion create feature/new-feature
  scion create feature/payment --base develop
  scion create bugfix/issue-123 --remote upstream
  scion create feature/refactor --force`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createBaseBranch, "base", "b", "", "ベースブランチを指定 (デフォルト: 設定ファイルの値または現在のブランチ)")
	createCmd.Flags().StringVarP(&createRemote, "remote", "r", "", "リモートリポジトリを指定 (デフォルト: origin)")
	createCmd.Flags().BoolVarP(&createForce, "force", "f", false, "既存のworktreeを強制的に上書き")
}

func runCreate(cmd *cobra.Command, args []string) error {
	branchName := args[0]

	// Gitリポジトリかどうか確認
	if !git.IsGitRepository() {
		return fmt.Errorf("Gitリポジトリ内で実行してください")
	}

	// リポジトリルートを取得
	repoRoot, err := git.GetRepositoryRoot()
	if err != nil {
		return err
	}

	// 設定からデフォルト値を取得
	config := GetConfig()
	baseDir := config.Worktree.BaseDir
	baseBranch := createBaseBranch
	remote := createRemote

	// フラグが指定されていない場合は設定ファイルの値を使用
	if baseBranch == "" {
		baseBranch = config.Git.DefaultBaseBranch
	}
	if remote == "" {
		remote = config.Git.DefaultRemote
	}

	// fetch を実行（設定で有効な場合）
	if config.Git.FetchBeforeCreate {
		output.Info("リモートから最新の情報を取得しています...")
		if err := git.Fetch(remote); err != nil {
			output.Warning("fetchに失敗しました: %v", err)
		}
	}

	// worktreeパスを構築
	// ベースリポジトリの親ディレクトリに wtree ディレクトリを作成
	parentDir := filepath.Dir(repoRoot)

	// ブランチ名からパスセーフな名前を生成
	safeBranchName := strings.ReplaceAll(branchName, "/", "-")
	worktreePath := filepath.Join(parentDir, baseDir, safeBranchName)

	// ブランチが既に存在するか確認
	branchExists := git.BranchExists(branchName)
	worktreeExists := git.WorktreeExists(worktreePath)

	if worktreeExists && !createForce {
		return fmt.Errorf("worktree '%s' は既に存在します\n--force オプションで上書きできます", worktreePath)
	}

	if branchExists && !createForce {
		output.Warning("ブランチ '%s' は既に存在します。既存のブランチをチェックアウトします", branchName)
	}

	// 強制モードで既存のworktreeがある場合は削除
	if worktreeExists && createForce {
		output.Info("既存のworktreeを削除しています...")
		if err := git.RemoveWorktree(worktreePath, true); err != nil {
			return fmt.Errorf("既存のworktreeの削除に失敗しました: %w", err)
		}
	}

	// worktree を作成
	output.Info("worktreeを作成しています: %s", branchName)
	if err := git.CreateWorktree(worktreePath, branchName, baseBranch, createForce); err != nil {
		return err
	}

	output.Success("Worktree ディレクトリを作成しました: %s", worktreePath)
	if !branchExists {
		output.Success("ブランチ '%s' を作成してチェックアウトしました", branchName)
	} else {
		output.Success("ブランチ '%s' をチェックアウトしました", branchName)
	}
	output.Info("パス: %s", worktreePath)

	return nil
}
