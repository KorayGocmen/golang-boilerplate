package context

import "testing"

func TestLang(t *testing.T) {
	cases := []struct {
		lang string
		want Language
	}{
		{"tr", LangTR},
		{"en", LangEN},
		{"", LangTR},
		{"foo", LangTR},
	}

	for _, c := range cases {
		ctx := WithValue(Background(), KeyLang, ToLanguage(c.lang))
		if got := Lang(ctx); got != c.want {
			t.Errorf("Lang(%q) == %q, want %q", c.lang, got, c.want)
		}
	}
}
