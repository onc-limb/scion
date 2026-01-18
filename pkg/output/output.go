package output

import (
	"fmt"
	"os"
)

// ColorEnabled はカラー出力が有効かどうか
var ColorEnabled = true

// ANSI カラーコード
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// Success は成功メッセージを出力する
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if ColorEnabled {
		fmt.Printf("%s✓ %s%s\n", colorGreen, msg, colorReset)
	} else {
		fmt.Printf("✓ %s\n", msg)
	}
}

// Error はエラーメッセージを出力する
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if ColorEnabled {
		fmt.Fprintf(os.Stderr, "%s✗ %s%s\n", colorRed, msg, colorReset)
	} else {
		fmt.Fprintf(os.Stderr, "✗ %s\n", msg)
	}
}

// Warning は警告メッセージを出力する
func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if ColorEnabled {
		fmt.Printf("%s⚠ %s%s\n", colorYellow, msg, colorReset)
	} else {
		fmt.Printf("⚠ %s\n", msg)
	}
}

// Info は情報メッセージを出力する
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if ColorEnabled {
		fmt.Printf("%s→ %s%s\n", colorCyan, msg, colorReset)
	} else {
		fmt.Printf("→ %s\n", msg)
	}
}

// Print は通常のメッセージを出力する
func Print(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// SetColorEnabled はカラー出力の有効/無効を設定する
func SetColorEnabled(enabled bool) {
	ColorEnabled = enabled
}
