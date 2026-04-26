package domain

import (
	"strings"
	"testing"
)

func TestValidateProcessRequest(t *testing.T) {
	longRunes := strings.Repeat("中", 25) // 75 bytes, 25 runes
	longBytes := strings.Repeat("a", MaxInputTextBytes+1)

	tests := []struct {
		name    string
		in      *ProcessRequest
		wantErr bool
	}{
		{
			name:    "nil",
			in:      nil,
			wantErr: true,
		},
		{
			name:    "empty",
			in:      &ProcessRequest{Text: ""},
			wantErr: true,
		},
		{
			name:    "whitespace",
			in:      &ProcessRequest{Text: "   \t  "},
			wantErr: true,
		},
		{
			name:    "too few runes",
			in:      &ProcessRequest{Text: "short text"},
			wantErr: true,
		},
		{
			name:    "too many bytes",
			in:      &ProcessRequest{Text: longBytes},
			wantErr: true,
		},
		{
			name:    "valid at min",
			in:      &ProcessRequest{Text: longRunes + "x"},
			wantErr: false,
		},
		{
			name:    "valid longer",
			in:      &ProcessRequest{Text: strings.Repeat("word ", 10)}, // 50+ runes
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProcessRequest(tc.in)
			if tc.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if err != nil {
				if _, ok := err.(*ValidationError); !ok {
					t.Fatalf("want *ValidationError, got %T", err)
				}
			}
		})
	}
}

func TestValidateProcessResponse(t *testing.T) {
	mk3 := func(a, b, c string) []string { return []string{a, b, c} }
	long := strings.Repeat("a", MaxSummaryRunes+1)
	longItem := strings.Repeat("b", MaxActionItemRunes+1)

	tests := []struct {
		name    string
		in      *ProcessResponse
		wantErr bool
	}{
		{
			name:    "nil",
			in:      nil,
			wantErr: true,
		},
		{
			name:    "no summary",
			in:      &ProcessResponse{Summary: "   ", ActionItems: mk3("1", "2", "3")},
			wantErr: true,
		},
		{
			name:    "summary too long",
			in:      &ProcessResponse{Summary: long, ActionItems: mk3("1", "2", "3")},
			wantErr: true,
		},
		{
			name:    "wrong number of action items 2",
			in:      &ProcessResponse{Summary: "ok", ActionItems: []string{"a", "b"}},
			wantErr: true,
		},
		{
			name:    "wrong number of action items 4",
			in:      &ProcessResponse{Summary: "ok", ActionItems: []string{"a", "b", "c", "d"}},
			wantErr: true,
		},
		{
			name:    "empty action",
			in:      &ProcessResponse{Summary: "ok", ActionItems: mk3("1", "  ", "3")},
			wantErr: true,
		},
		{
			name:    "action too long",
			in:      &ProcessResponse{Summary: "ok", ActionItems: mk3("1", "2", longItem)},
			wantErr: true,
		},
		{
			name: "valid",
			in: &ProcessResponse{
				Summary:     "A reasonable summary of the work.",
				ActionItems: mk3("First task", "Second task", "Third task"),
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProcessResponse(tc.in)
			if tc.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if err != nil {
				if _, ok := err.(*ValidationError); !ok {
					t.Fatalf("want *ValidationError, got %T", err)
				}
			}
		})
	}
}
