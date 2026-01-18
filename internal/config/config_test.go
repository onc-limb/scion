package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Worktree.BaseDir != "wtree" {
		t.Errorf("expected BaseDir to be 'wtree', got '%s'", cfg.Worktree.BaseDir)
	}

	if cfg.Git.DefaultRemote != "origin" {
		t.Errorf("expected DefaultRemote to be 'origin', got '%s'", cfg.Git.DefaultRemote)
	}

	if cfg.Git.DefaultBaseBranch != "main" {
		t.Errorf("expected DefaultBaseBranch to be 'main', got '%s'", cfg.Git.DefaultBaseBranch)
	}

	if !cfg.Worktree.AutoCreateDir {
		t.Error("expected AutoCreateDir to be true")
	}

	if !cfg.UI.ColorOutput {
		t.Error("expected ColorOutput to be true")
	}
}

func TestSaveAndLoad(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	// 設定を作成して保存
	cfg := DefaultConfig()
	cfg.Worktree.BaseDir = "custom-wtree"
	cfg.Git.DefaultRemote = "upstream"

	if err := Save(cfg, configPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// ファイルが作成されたことを確認
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	// 設定を読み込み
	loadedCfg := DefaultConfig()
	if err := loadFromFile(configPath, loadedCfg); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// 値が正しく読み込まれたことを確認
	if loadedCfg.Worktree.BaseDir != "custom-wtree" {
		t.Errorf("expected BaseDir to be 'custom-wtree', got '%s'", loadedCfg.Worktree.BaseDir)
	}

	if loadedCfg.Git.DefaultRemote != "upstream" {
		t.Errorf("expected DefaultRemote to be 'upstream', got '%s'", loadedCfg.Git.DefaultRemote)
	}
}

func TestGlobalConfigPath(t *testing.T) {
	path, err := GlobalConfigPath()
	if err != nil {
		t.Fatalf("failed to get global config path: %v", err)
	}

	// パスが ~/.config/scion/config.toml で終わることを確認
	if filepath.Base(path) != "config.toml" {
		t.Errorf("expected config file name to be 'config.toml', got '%s'", filepath.Base(path))
	}

	if filepath.Base(filepath.Dir(path)) != "scion" {
		t.Errorf("expected parent directory to be 'scion', got '%s'", filepath.Base(filepath.Dir(path)))
	}
}

func TestLocalConfigPath(t *testing.T) {
	path := LocalConfigPath()

	if path != filepath.Join(".scion", "config.toml") {
		t.Errorf("expected local config path to be '.scion/config.toml', got '%s'", path)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	cfg := DefaultConfig()
	err := loadFromFile("/non/existent/path/config.toml", cfg)

	if err == nil {
		t.Error("expected error when loading non-existent file")
	}

	if !os.IsNotExist(err) {
		t.Errorf("expected os.IsNotExist error, got: %v", err)
	}
}
