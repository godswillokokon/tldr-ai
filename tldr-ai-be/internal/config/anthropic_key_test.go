package config

import "testing"

func TestIsAnthropicPlaceholderKey(t *testing.T) {
	if !IsAnthropicPlaceholderKey("") {
		t.Fatal("empty should be placeholder")
	}
	if !IsAnthropicPlaceholderKey("  example-key-please-replace  ") {
		t.Fatal("example keyword")
	}
	if !IsAnthropicPlaceholderKey("sk-ant-xxxxxxxxxxxxxxxxxxxxxxxx") {
		t.Fatal("all-mask sk-ant- tail should be placeholder")
	}
}
