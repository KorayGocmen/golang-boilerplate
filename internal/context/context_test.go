package context

import (
	"testing"
)

func TestBackground(t *testing.T) {
	ctx := Background()
	if ctx == nil {
		t.Fatalf("want: Background() != nil; got: Background returned nil")
	}
}

func TestWithValue(t *testing.T) {
	ctx := WithValue(Background(), "key", "value")
	if ctx == nil {
		t.Fatalf("want: WithValue() != nil; got: WithValue returned nil")
	}

	if got := ctx.Value("key"); got != "value" {
		t.Fatalf(`want: ctx.Value("key") = "value"; got "%s"`, got)
	}
}

func TestWithCancel(t *testing.T) {
	ctx, cancel := WithCancel(Background())
	if ctx == nil {
		t.Fatalf("want: WithCancel() != nil; got: WithCancel returned nil")
	}
	cancel()
}
