package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/ongasatoshi/scion/internal/config"
	"github.com/ongasatoshi/scion/pkg/output"
	"github.com/spf13/cobra"
)

var (
	configGlobal bool
	configLocal  bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "scionの設定を管理",
	Long: `config コマンドはscionの設定を管理します。

サブコマンド:
  get <key>      - 特定の設定値を取得
  set <key> <value> - 設定値を更新
  list           - すべての設定を表示
  reset          - 設定をデフォルトに戻す
  edit           - エディタで設定ファイルを開く

例:
  scion config get worktree.base_dir
  scion config set ui.color_output false
  scion config list
  scion config edit`,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "特定の設定値を取得",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "設定値を更新",
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "すべての設定を表示",
	Args:  cobra.NoArgs,
	RunE:  runConfigList,
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "設定をデフォルトに戻す",
	Args:  cobra.NoArgs,
	RunE:  runConfigReset,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "エディタで設定ファイルを開く",
	Args:  cobra.NoArgs,
	RunE:  runConfigEdit,
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().BoolVar(&configGlobal, "global", false, "グローバル設定を対象とする")
	configCmd.PersistentFlags().BoolVar(&configLocal, "local", false, "ローカル設定を対象とする")

	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configEditCmd)
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]
	cfg := GetConfig()

	value, err := getConfigValue(cfg, key)
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// 設定ファイルのパスを決定
	var configPath string
	var err error

	if configLocal {
		configPath = config.LocalConfigPath()
	} else {
		configPath, err = config.GlobalConfigPath()
		if err != nil {
			return err
		}
	}

	// 現在の設定を読み込む
	cfg, err := config.Load("")
	if err != nil {
		return err
	}

	// 設定値を更新
	if err := setConfigValue(cfg, key, value); err != nil {
		return err
	}

	// 設定を保存
	if err := config.Save(cfg, configPath); err != nil {
		return err
	}

	output.Success("設定を更新しました: %s = %s", key, value)
	return nil
}

func runConfigList(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// グローバル設定を表示
	globalPath, err := config.GlobalConfigPath()
	if err == nil {
		fmt.Printf("グローバル設定 (%s):\n", globalPath)
		printConfigSection(cfg, "  ")
	}

	// ローカル設定の存在を確認して表示
	localPath := config.LocalConfigPath()
	if _, err := os.Stat(localPath); err == nil {
		fmt.Printf("\nローカル設定 (%s):\n", localPath)
		fmt.Println("  (ローカル設定が存在します)")
	}

	return nil
}

func runConfigReset(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// 確認プロンプト
	if cfg.UI.ConfirmDestructive {
		fmt.Print("設定をデフォルトに戻しますか? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			output.Info("キャンセルしました")
			return nil
		}
	}

	// デフォルト設定を取得
	defaultCfg := config.DefaultConfig()

	// 設定ファイルのパスを決定
	var configPath string
	var err error

	if configLocal {
		configPath = config.LocalConfigPath()
	} else {
		configPath, err = config.GlobalConfigPath()
		if err != nil {
			return err
		}
	}

	// デフォルト設定を保存
	if err := config.Save(defaultCfg, configPath); err != nil {
		return err
	}

	output.Success("設定をデフォルトにリセットしました")
	return nil
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	// エディタを決定
	editor := os.Getenv("EDITOR")
	if editor == "" {
		cfg := GetConfig()
		editor = cfg.Editor.Command
	}

	// 設定ファイルのパスを決定
	var configPath string
	var err error

	if configLocal {
		configPath = config.LocalConfigPath()
	} else {
		configPath, err = config.GlobalConfigPath()
		if err != nil {
			return err
		}
	}

	// 設定ファイルが存在しない場合は作成
	if err := config.EnsureConfigExists(); err != nil {
		return err
	}

	// エディタを起動
	execCmd := exec.Command(editor, configPath)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}

func getConfigValue(cfg *config.Config, key string) (string, error) {
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("無効なキー形式です: %s (例: worktree.base_dir)", key)
	}

	section := parts[0]
	field := parts[1]

	v := reflect.ValueOf(cfg).Elem()

	// セクションを見つける
	sectionField := v.FieldByNameFunc(func(name string) bool {
		return strings.EqualFold(name, section)
	})
	if !sectionField.IsValid() {
		return "", fmt.Errorf("セクション '%s' が見つかりません", section)
	}

	// フィールドを見つける
	fieldValue := sectionField.FieldByNameFunc(func(name string) bool {
		return strings.EqualFold(name, field) || strings.EqualFold(toSnakeCase(name), field)
	})
	if !fieldValue.IsValid() {
		return "", fmt.Errorf("フィールド '%s' が見つかりません", field)
	}

	return fmt.Sprintf("%v", fieldValue.Interface()), nil
}

func setConfigValue(cfg *config.Config, key, value string) error {
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return fmt.Errorf("無効なキー形式です: %s (例: worktree.base_dir)", key)
	}

	section := parts[0]
	field := parts[1]

	v := reflect.ValueOf(cfg).Elem()

	// セクションを見つける
	sectionField := v.FieldByNameFunc(func(name string) bool {
		return strings.EqualFold(name, section)
	})
	if !sectionField.IsValid() {
		return fmt.Errorf("セクション '%s' が見つかりません", section)
	}

	// フィールドを見つける
	fieldValue := sectionField.FieldByNameFunc(func(name string) bool {
		return strings.EqualFold(name, field) || strings.EqualFold(toSnakeCase(name), field)
	})
	if !fieldValue.IsValid() {
		return fmt.Errorf("フィールド '%s' が見つかりません", field)
	}

	// 値を設定
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("無効なbool値です: %s", value)
		}
		fieldValue.SetBool(boolValue)
	case reflect.Int, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("無効な数値です: %s", value)
		}
		fieldValue.SetInt(intValue)
	default:
		return fmt.Errorf("サポートされていない型です: %v", fieldValue.Kind())
	}

	return nil
}

func printConfigSection(cfg *config.Config, indent string) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		fmt.Printf("%s%s:\n", indent, strings.ToLower(fieldType.Name))

		if field.Kind() == reflect.Struct {
			printStructFields(field, indent+"  ")
		}
	}
}

func printStructFields(v reflect.Value, indent string) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("toml")
		if tag == "" {
			tag = toSnakeCase(fieldType.Name)
		}

		fmt.Printf("%s%s: %v\n", indent, tag, field.Interface())
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
