package duration

import (
	"testing"
	"time"
)

func TestSeconds(t *testing.T) {
	if got := Seconds(1); got != 1*time.Second {
		t.Fatalf("Seconds(1) = %d; want 1", got)
	}
}
