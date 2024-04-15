package context

import (
	"testing"
)

func TestMap(t *testing.T) {
	ctx := Background()
	ctx = WithValue(ctx, KeyLogLevel, 1)
	ctx = WithValue(ctx, KeyMethod, "method")
	ctx = WithValue(ctx, KeyPath, "path")
	ctx = WithValue(ctx, KeyRequestID, "request_id")
	ctx = WithValue(ctx, KeyReqBody, nil)

	got := Map(ctx)
	want := map[ContextKey]interface{}{
		"log_level":  1,
		"method":     "method",
		"path":       "path",
		"request_id": "request_id",
	}

	if len(got) != len(want) {
		t.Fatalf("want: len(want) = %d; got: %d", len(want), len(got))
	}

	for k, v := range want {
		if got[k] != v {
			t.Fatalf(`want: got[%q] = %q; got: %q`, k, v, got[k])
		}
	}
}

func TestString(t *testing.T) {
	ctx := Background()
	ctx = WithValue(ctx, KeyLogLevel, 1)
	ctx = WithValue(ctx, KeyMethod, "method")
	ctx = WithValue(ctx, KeyPath, "path")
	ctx = WithValue(ctx, KeyRequestID, "request_id")

	want := `log_level="1" method="method" path="path" request_id="request_id" format`
	if got := String(ctx, "format"); got != want {
		t.Fatalf("want: %s; got: %s", want, got)
	}
}
