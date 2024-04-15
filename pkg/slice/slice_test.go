package slice

import "testing"

func TestContains(t *testing.T) {
	if !Contains("a", []string{"a", "b", "c"}) {
		t.Fatalf("want: true, got: false")
	}

	if Contains("d", []string{"a", "b", "c"}) {
		t.Fatalf("want: false, got: true")
	}

	if !Contains(1, []int{1, 2, 3}) {
		t.Fatalf("want: true, got: false")
	}

	if Contains(4, []int{1, 2, 3}) {
		t.Fatalf("want: false, got: true")
	}

	if !Contains(int64(1), []int64{1, 2, 3}) {
		t.Fatalf("want: true, got: false")
	}

	if Contains(int64(4), []int64{1, 2, 3}) {
		t.Fatalf("want: false, got: true")
	}

	if !Contains(float64(1), []float64{1, 2, 3}) {
		t.Fatalf("want: true, got: false")
	}

	if Contains(float64(4), []float64{1, 2, 3}) {
		t.Fatalf("want: false, got: true")
	}
}

func TestMin(t *testing.T) {
	if Min(1) != 1 {
		t.Fatalf("want: 1, got: %d", Min(1))
	}

	if Min(1, 2, 3) != 1 {
		t.Fatalf("want: 1, got: %d", Min(1, 2, 3))
	}

	if Min(int64(1), int64(2), int64(3)) != int64(1) {
		t.Fatalf("want: 1, got: %d", Min(int64(1), int64(2), int64(3)))
	}

	if Min(float64(1), float64(2), float64(3)) != float64(1) {
		t.Fatalf("want: 1, got: %f", Min(float64(1), float64(2), float64(3)))
	}
}
