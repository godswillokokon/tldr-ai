package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"tldr-ai-be/internal/domain"
	"tldr-ai-be/internal/errs"
)

type fakeProvider struct {
	out      string
	err      error
	modelTag string
}

func (f *fakeProvider) Complete(ctx context.Context, prompt string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if f.out != "" {
		return f.out, nil
	}
	return "{}", nil
}

func (f *fakeProvider) ModelTag() string {
	if f.modelTag != "" {
		return f.modelTag
	}
	return "test-model"
}

func longInput() *domain.ProcessRequest {
	return &domain.ProcessRequest{Text: "This is long enough input text for the minimum rune count rule here."}
}

func validModelJSON() string {
	return `{"summary":"A concise summary in two short sentences. It describes the work done clearly and helpfully.","actionItems":["Do the first follow up task in plain language without numbering","Do the second follow up in plain language without bullets","Do the third follow up in plain language without a prefix"]}`
}

func TestTextProcessor_OK(t *testing.T) {
	tp := NewTextProcessor(&fakeProvider{out: validModelJSON(), modelTag: "m1"})
	res, err := tp.Process(context.Background(), longInput())
	if err != nil {
		t.Fatal(err)
	}
	if res.Model != "m1" {
		t.Fatalf("Model: %q", res.Model)
	}
	if len(res.ActionItems) != 3 {
		t.Fatalf("actions: %d", len(res.ActionItems))
	}
}

func TestTextProcessor_validateRequest(t *testing.T) {
	tp := NewTextProcessor(&fakeProvider{out: validModelJSON()})
	_, err := tp.Process(context.Background(), &domain.ProcessRequest{Text: "short"})
	if err == nil {
		t.Fatal("expected error")
	}
	var ve *domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("got %T", err)
	}
}

func TestTextProcessor_completeErr(t *testing.T) {
	tp := NewTextProcessor(&fakeProvider{err: &errs.HTTPError{PublicMessage: "upstream", Status: http.StatusBadGateway}})
	_, err := tp.Process(context.Background(), longInput())
	if err == nil {
		t.Fatal("expected error")
	}
	var he *errs.HTTPError
	if !errors.As(err, &he) || he.PublicMessage != "upstream" || he.Status != http.StatusBadGateway {
		t.Fatalf("err: %v", err)
	}
}

func TestTextProcessor_parseErr(t *testing.T) {
	tp := NewTextProcessor(&fakeProvider{out: "not valid json at all"})
	_, err := tp.Process(context.Background(), longInput())
	if err == nil {
		t.Fatal("expected error")
	}
	var aio *domain.InvalidAIOutputError
	if !errors.As(err, &aio) {
		t.Fatalf("want InvalidAIOutputError, got %T: %v", err, err)
	}
}

func TestTextProcessor_validateResponse(t *testing.T) {
	// JSON parses but not exactly 3 items after domain rules.
	bad := `{"summary":"A concise summary in two short sentences. It describes the work done.","actionItems":["one","two"]}`
	tp := NewTextProcessor(&fakeProvider{out: bad})
	_, err := tp.Process(context.Background(), longInput())
	if err == nil {
		t.Fatal("expected error")
	}
	var ve *domain.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("want ValidationError, got %T: %v", err, err)
	}
}

func TestTextProcessor_modelTagDefault(t *testing.T) {
	tp := NewTextProcessor(&fakeProvider{out: validModelJSON()})
	res, err := tp.Process(context.Background(), longInput())
	if err != nil {
		t.Fatal(err)
	}
	if res.Model != "test-model" {
		t.Fatalf("Model: %q", res.Model)
	}
}
