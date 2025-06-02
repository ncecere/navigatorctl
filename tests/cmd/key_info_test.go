// tests/cmd/key_info_test.go

package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestKeyInfo_Expected(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "key", "info", "--key", "58c9b3049f941b9e40d35ae045ff47040ff47d521b6aad3b70bc1e6ddccf150b", "--api-url", "https://ai.bitop.dev", "--api-key", "sk-6425")
	cmd.Dir = "../"
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error, got %v, output: %s", err, out.String())
	}
	if !strings.Contains(out.String(), "Key Info:") {
		t.Errorf("Expected output to contain 'Key Info:', got: %s", out.String())
	}
}

func TestKeyInfo_MissingKey(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "key", "info", "--api-url", "https://ai.bitop.dev", "--api-key", "sk-6425")
	cmd.Dir = "../"
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err == nil {
		t.Fatalf("Expected error due to missing --key, got none")
	}
	if !strings.Contains(out.String(), "API URL, API Key, and --key are required") {
		t.Errorf("Expected error message for missing --key, got: %s", out.String())
	}
}

func TestKeyInfo_InvalidKey(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "key", "info", "--key", "invalidkey", "--api-url", "https://ai.bitop.dev", "--api-key", "sk-6425")
	cmd.Dir = "../"
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err == nil {
		t.Fatalf("Expected error due to invalid key, got none")
	}
	if !strings.Contains(out.String(), "API error:") {
		t.Errorf("Expected error message for invalid key, got: %s", out.String())
	}
}
