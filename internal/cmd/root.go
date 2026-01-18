package cmd

import (
	"fmt"

	"github.com/ongasatoshi/scion/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Version はビルド時に設定される
	Version = "0.1.0"
	// Commit はビルド時に設定される
	Commit = "unknown"

	cfgFile string
	cfg     *config.Config
)

// rootCmd はベースコマンド
var rootCmd = &cobra.Command{
	Use:   "scion",
	Short: "Git Worktreeを効率的に管理するCLIツール",
	Long: `scionはGit Worktreeを効率的に管理するためのCLIツールです。
AI開発環境（Claude Code、GitHub Copilot、Cursor）との連携を前提とした、
シンプルで使いやすいワークフロー管理を提供します。

利用可能なコマンド:
  create  - 新しいworktreeブランチを作成
  clear   - 既存のworktreeブランチを削除
  config  - scionの設定を管理`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
		}
		return nil
	},
}

// Execute はルートコマンドを実行する
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "設定ファイルのパス (デフォルト: ~/.config/scion/config.toml)")
	rootCmd.Flags().BoolP("version", "v", false, "バージョン情報を表示")

	rootCmd.SetVersionTemplate(fmt.Sprintf("scion version %s (commit: %s)\n", Version, Commit))
	rootCmd.Version = Version
}

// GetConfig は現在の設定を返す
func GetConfig() *config.Config {
	return cfg
}
