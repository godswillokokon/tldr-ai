package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tldr-ai-be/internal/domain"
	"tldr-ai-be/internal/errs"
)

func TestDecodeProcessRequest(t *testing.T) {
	good := `{"text":"short text that is long enough"}`

	tests := []struct {
		name     string
		body     string
		maxBytes int64
		want     *domain.ProcessRequest
		wantCode int
		wantErr  string // public message of *errs.HTTPError, or "ok"
	}{
		{
			name:     "ok",
			body:     good,
			maxBytes: 1024,
			want:     &domain.ProcessRequest{Text: "short text that is long enough"},
			wantCode: 0,
			wantErr:  "ok",
		},
		{
			name:     "empty body",
			body:     "",
			maxBytes: 100,
			wantCode: 400,
			wantErr:  "Request body is empty",
		},
		{
			name:     "too large",
			body:     "1234567",
			maxBytes: 5,
			wantCode: 413,
			wantErr:  "Request body is too large",
		},
		{
			name:     "unknown field",
			body:     `{"text":"short text that is long enough","extra":true}`,
			maxBytes: 200,
			wantCode: 400,
			wantErr:  "Invalid JSON in request body",
		},
		{
			name:     "trailing value",
			body:     `{"text":"short text that is long enough"}{"x":1}`,
			maxBytes: 200,
			wantCode: 400,
			wantErr:  "Request body contains additional JSON data",
		},
		{
			name:     "trailing with whitespace",
			body:     `{"text":"short text that is long enough"}  false`,
			maxBytes: 200,
			wantCode: 400,
			wantErr:  "Request body contains additional JSON data",
		},
		{
			name:     "garbage not json",
			body:     `notjson`,
			maxBytes: 100,
			wantCode: 400,
			wantErr:  "Invalid JSON in request body",
		},
		{
			name:     "wrong type for text",
			body:     `{"text":true}`,
			maxBytes: 100,
			wantCode: 400,
			wantErr:  "Invalid JSON in request body",
		},
		{
			name:     "at exact max size",
			body:     "012345",
			maxBytes: 6,
			wantCode: 400, // not valid json but we check size only after read — invalid
			wantErr:  "Invalid JSON in request body",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/process", strings.NewReader(tc.body))
			req.ContentLength = int64(len(tc.body))
			out, err := DecodeProcessRequest(req, tc.maxBytes)
			if tc.wantErr == "ok" {
				if err != nil {
					t.Fatalf("err: %v", err)
				}
				if out.Text != tc.want.Text {
					t.Fatalf("Text: %q, want %q", out.Text, tc.want.Text)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error")
			}
			var he *errs.HTTPError
			if !errors.As(err, &he) {
				t.Fatalf("want *errs.HTTPError, got %T: %v", err, err)
			}
			if he.Status != tc.wantCode {
				t.Fatalf("status: %d, want %d (err=%q)", he.Status, tc.wantCode, he.PublicMessage)
			}
			if he.PublicMessage != tc.wantErr {
				t.Fatalf("message: %q, want %q", he.PublicMessage, tc.wantErr)
			}
		})
	}
}
