// tests/cmd/key_list_test.go

package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestKeyList_Expected(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "key", "list", "--api-url", "https://ai.bitop.dev", "--api-key", "sk-6425")
	cmd.Dir = "../"
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error, got %v, output: %s", err, out.String())
	}
	if !strings.Contains(out.String(), "Keys:") {
		t.Errorf("Expected output to contain 'Keys:', got: %s", out.String())
	}
}

func TestKeyList_MissingAPIKey(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "key", "list", "--api-url", "https://ai.bitop.dev")
	cmd.Dir = "../"
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err == nil {
		t.Fatalf("Expected error due to missing API key, got none")
	}
	if !strings.Contains(out.String(), "API URL and API Key are required") {
		t.Errorf("Expected error message for missing API key, got: %s", out.String())
	}
}

func TestKeyList_InvalidURL(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "key", "list", "--api-url", "http://invalid-url", "--api-key", "sk-6425")
	cmd.Dir = "../"
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err == nil {
		t.Fatalf("Expected error due to invalid URL, got none")
	}
	if !strings.Contains(out.String(), "Request failed:") && !strings.Contains(out.String(), "API error:") {
		t.Errorf("Expected error message for invalid URL, got: %s", out.String())
	}
}
