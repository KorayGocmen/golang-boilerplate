package array

import "testing"

func TestFromTo(t *testing.T) {
	list := BeginEndInclude(1, 10)
	if len(list) != 10 {
		t.Fatalf("list length mismatch")
	}

	if list[0] != 1 {
		t.Fatalf("list[0] mismatch")
	}

	if list[9] != 10 {
		t.Fatalf("list[9] mismatch")
	}
}
