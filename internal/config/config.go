package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config はscionの設定構造体
type Config struct {
	Repository RepositoryConfig `toml:"repository"`
	Worktree   WorktreeConfig   `toml:"worktree"`
	Git        GitConfig        `toml:"git"`
	UI         UIConfig         `toml:"ui"`
	Editor     EditorConfig     `toml:"editor"`
}

// RepositoryConfig はリポジトリ関連の設定
type RepositoryConfig struct {
	BaseRepository string `toml:"base_repository"`
	BaseBranch     string `toml:"base_branch"`
}

// WorktreeConfig はworktree関連の設定
type WorktreeConfig struct {
	BaseDir                string `toml:"base_dir"`
	AutoCreateDir          bool   `toml:"auto_create_dir"`
	CleanupOnBranchDelete  bool   `toml:"cleanup_on_branch_delete"`
}

// GitConfig はGit関連の設定
type GitConfig struct {
	DefaultRemote     string `toml:"default_remote"`
	DefaultBaseBranch string `toml:"default_base_branch"`
	FetchBeforeCreate bool   `toml:"fetch_before_create"`
}

// UIConfig はUI関連の設定
type UIConfig struct {
	ColorOutput        bool `toml:"color_output"`
	Verbose            bool `toml:"verbose"`
	ConfirmDestructive bool `toml:"confirm_destructive"`
}

// EditorConfig はエディタ設定
type EditorConfig struct {
	Command string `toml:"command"`
}

// DefaultConfig はデフォルト設定を返す
func DefaultConfig() *Config {
	return &Config{
		Repository: RepositoryConfig{
			BaseRepository: "",
			BaseBranch:     "main",
		},
		Worktree: WorktreeConfig{
			BaseDir:               "wtree",
			AutoCreateDir:         true,
			CleanupOnBranchDelete: true,
		},
		Git: GitConfig{
			DefaultRemote:     "origin",
			DefaultBaseBranch: "main",
			FetchBeforeCreate: true,
		},
		UI: UIConfig{
			ColorOutput:        true,
			Verbose:            false,
			ConfirmDestructive: true,
		},
		Editor: EditorConfig{
			Command: "vi",
		},
	}
}

// GlobalConfigPath はグローバル設定ファイルのパスを返す
func GlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "scion", "config.toml"), nil
}

// LocalConfigPath はローカル設定ファイルのパスを返す
func LocalConfigPath() string {
	return filepath.Join(".scion", "config.toml")
}

// Load は設定ファイルを読み込む
func Load(customPath string) (*Config, error) {
	cfg := DefaultConfig()

	// グローバル設定の読み込み
	globalPath, err := GlobalConfigPath()
	if err == nil {
		if err := loadFromFile(globalPath, cfg); err != nil && !os.IsNotExist(err) {
			return cfg, err
		}
	}

	// ローカル設定の読み込み（上書き）
	localPath := LocalConfigPath()
	if err := loadFromFile(localPath, cfg); err != nil && !os.IsNotExist(err) {
		return cfg, err
	}

	// カスタムパスが指定されている場合
	if customPath != "" {
		if err := loadFromFile(customPath, cfg); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}

// loadFromFile はファイルから設定を読み込む
func loadFromFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, cfg)
}

// Save は設定をファイルに保存する
func Save(cfg *Config, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// EnsureConfigExists は設定ファイルが存在しない場合に作成する
func EnsureConfigExists() error {
	path, err := GlobalConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := DefaultConfig()
		return Save(cfg, path)
	}

	return nil
}
