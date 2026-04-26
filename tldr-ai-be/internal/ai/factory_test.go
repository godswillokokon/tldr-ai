package ai

import "testing"

func TestNewProviderFromEnv(t *testing.T) {
	_, err := NewProviderFromEnv(func(string) string { return "" })
	if err == nil {
		t.Fatal("expected error for missing key")
	}

	_, err = NewProviderFromEnv(func(string) string { return "  placeholder-example  " })
	if err == nil {
		t.Fatal("expected error for placeholder key")
	}

	p, err := NewProviderFromEnv(func(k string) string {
		switch k {
		case "ANTHROPIC_API_KEY":
			return "sk-ant-api03-aaaaaaaaaaaaaaaaaaaaaa"
		case "ANTHROPIC_MODEL":
			return "my-model"
		default:
			return ""
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.ModelTag() != "my-model" {
		t.Fatalf("model: %q", p.ModelTag())
	}
}
