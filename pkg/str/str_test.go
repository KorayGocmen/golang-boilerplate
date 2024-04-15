package str

import "testing"

func TestSnakeToCamel(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{
			s:    "",
			want: "",
		},
		{
			s:    "hello",
			want: "Hello",
		},
		{
			s:    "hello_world",
			want: "HelloWorld",
		},
		{
			s:    "hello_world_123",
			want: "HelloWorld123",
		},
	}
	for _, tt := range tests {
		if got := SnakeToCamel(tt.s); got != tt.want {
			t.Fatalf("SnakeToCamel() = %v, want %v", got, tt.want)
		}
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		alphabet string
		s1       string
		s2       string
		want     int
	}{
		{
			alphabet: "abcdefghijklmnopqrstuvwxyz",
			s1:       "hello",
			s2:       "world",
			want:     -1,
		},
		{
			alphabet: "abcdefghijklmnopqrstuvwxyz",
			s1:       "world",
			s2:       "hello",
			want:     1,
		},
		{
			alphabet: "abcdefghijklmnopqrstuvwxyz",
			s1:       "hello",
			s2:       "hello",
			want:     0,
		},
		{
			alphabet: "abcdefghijklmnopqrstuvwxyz",
			s1:       "hello",
			s2:       "hello_world",
			want:     -1,
		},
		{
			alphabet: "abcdefghijklmnopqrstuvwxyz",
			s1:       "hello_world",
			s2:       "hello",
			want:     1,
		},
		{
			alphabet: "abcçdefgğhıijklmnoöprsştuüvyz",
			s1:       "can",
			s2:       "çan",
			want:     -1,
		},
		{
			alphabet: "abcçdefgğhıijklmnoöprsştuüvyz",
			s1:       "çan",
			s2:       "can",
			want:     1,
		},
	}

	for _, c := range cases {
		if got := Compare(c.alphabet, c.s1, c.s2); got != c.want {
			t.Fatalf("Compare(%s, %s, %s) = %v, want %v", c.alphabet, c.s1, c.s2, got, c.want)
		}
	}
}
