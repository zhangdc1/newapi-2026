package service

import "testing"

func TestCountDeepSeekTextTokenOfficialHelloSample(t *testing.T) {
	tokens, ok := CountDeepSeekTextToken("Hello!")
	if !ok {
		t.Fatal("deepseek tokenizer failed to initialize")
	}
	if tokens != 2 {
		t.Fatalf("expected official DeepSeek sample Hello! to be 2 tokens, got %d", tokens)
	}
}

func TestCountDeepSeekTextTokenMixedSamples(t *testing.T) {
	tests := map[string]int{
		"The quick brown fox jumps over 13 lazy dogs.": 11,
		"你好，世界":        3,
		"abc12345!?\n": 5,
	}
	for text, want := range tests {
		got, ok := CountDeepSeekTextToken(text)
		if !ok {
			t.Fatal("deepseek tokenizer failed to initialize")
		}
		if got != want {
			t.Fatalf("CountDeepSeekTextToken(%q) = %d, want %d", text, got, want)
		}
	}
}

func TestDeepSeekTokenCounterOnlyMatchesDeepSeekModels(t *testing.T) {
	tests := map[string]bool{
		"deepseek-chat":                         true,
		"deepseek-reasoner":                     true,
		"deepseek-v4-flash":                     true,
		"deepseek-ai/DeepSeek-V3.1":             true,
		"accounts/fireworks/models/deepseek-r1": true,
		"gpt-4o":                                false,
		"claude-3-5-sonnet":                     false,
	}
	for model, want := range tests {
		if got := IsDeepSeekModel(model); got != want {
			t.Fatalf("IsDeepSeekModel(%q) = %v, want %v", model, got, want)
		}
	}
}

func TestCountTextTokenUsesDeepSeekTokenizer(t *testing.T) {
	if got := CountTextToken("Hello!", "deepseek-chat"); got != 2 {
		t.Fatalf("CountTextToken(Hello!, deepseek-chat) = %d, want 2", got)
	}
}
