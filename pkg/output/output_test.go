package output

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestSetColorEnabled(t *testing.T) {
	// 初期状態を保存
	original := ColorEnabled

	SetColorEnabled(false)
	if ColorEnabled {
		t.Error("expected ColorEnabled to be false")
	}

	SetColorEnabled(true)
	if !ColorEnabled {
		t.Error("expected ColorEnabled to be true")
	}

	// 元の状態に戻す
	ColorEnabled = original
}

func TestSuccessOutput(t *testing.T) {
	// カラーを無効にしてテスト
	original := ColorEnabled
	ColorEnabled = false
	defer func() { ColorEnabled = original }()

	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Success("test message")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got '%s'", output)
	}

	if !strings.Contains(output, "✓") {
		t.Errorf("expected output to contain '✓', got '%s'", output)
	}
}

func TestInfoOutput(t *testing.T) {
	original := ColorEnabled
	ColorEnabled = false
	defer func() { ColorEnabled = original }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Info("info message")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "info message") {
		t.Errorf("expected output to contain 'info message', got '%s'", output)
	}

	if !strings.Contains(output, "→") {
		t.Errorf("expected output to contain '→', got '%s'", output)
	}
}

func TestWarningOutput(t *testing.T) {
	original := ColorEnabled
	ColorEnabled = false
	defer func() { ColorEnabled = original }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Warning("warning message")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "warning message") {
		t.Errorf("expected output to contain 'warning message', got '%s'", output)
	}

	if !strings.Contains(output, "⚠") {
		t.Errorf("expected output to contain '⚠', got '%s'", output)
	}
}

func TestErrorOutput(t *testing.T) {
	original := ColorEnabled
	ColorEnabled = false
	defer func() { ColorEnabled = original }()

	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Error("error message")

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "error message") {
		t.Errorf("expected output to contain 'error message', got '%s'", output)
	}

	if !strings.Contains(output, "✗") {
		t.Errorf("expected output to contain '✗', got '%s'", output)
	}
}

func TestPrintOutput(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print("plain message")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "plain message") {
		t.Errorf("expected output to contain 'plain message', got '%s'", output)
	}
}

func TestFormatArgs(t *testing.T) {
	original := ColorEnabled
	ColorEnabled = false
	defer func() { ColorEnabled = original }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Success("message with %s and %d", "string", 42)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "message with string and 42") {
		t.Errorf("expected formatted output, got '%s'", output)
	}
}
