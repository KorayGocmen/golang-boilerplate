package date

import (
	"testing"
	"time"
)

func TestEqual(t *testing.T) {
	cases := []struct {
		date1 time.Time
		date2 time.Time
		want  bool
	}{
		{
			date1: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			date2: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			want:  true,
		},
		{
			date1: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			date2: time.Date(2018, 1, 2, 0, 0, 0, 0, time.UTC),
			want:  false,
		},
		{
			date1: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			date2: time.Date(2018, 1, 1, 4, 0, 0, 0, time.Local),
			want:  true,
		},
	}

	for _, c := range cases {
		if got := Equal(c.date1, c.date2); got != c.want {
			t.Fatalf("Equal(%v, %v) = %v; want %v", c.date1, c.date2, got, c.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	cases := []struct {
		date time.Time
		want time.Time
	}{
		{
			date: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			want: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			date: time.Date(2018, 1, 1, 4, 0, 0, 0, time.Local),
			want: time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local),
		},
	}

	for _, c := range cases {
		if got := Truncate(c.date); got != c.want {
			t.Fatalf("Truncate(%v) = %v; want %v", c.date, got, c.want)
		}
	}
}
