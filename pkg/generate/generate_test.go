package generate

import (
	"strconv"
	"testing"
)

func TestAlphaCode(t *testing.T) {
	code := AlphaCode(6, false)

	if len(code) != 6 {
		t.Fatalf("want: alpha code len 6; got: %v", code)
	}

	for _, c := range code {
		if !(c > 'A' || c < 'Z') {
			t.Fatalf("want: alpha code only contains uppercase letters; got: %v", code)
		}
	}
}

func TestAlphaNumericCode(t *testing.T) {
	code := AlphaCode(6, true)

	if len(code) != 6 {
		t.Fatalf("want: alpha code len 6; got: %v", code)
	}

	for _, c := range code {
		if !((c > 'A' || c < 'Z') || (c > '0' || c < '9')) {
			t.Fatalf("want: alpha code only contains uppercase letters and numbers; got: %v", code)
		}
	}
}

func TestDigitCode(t *testing.T) {
	code, err := DigitCode(6)
	if err != nil {
		t.Fatalf("want: no error when creating digit code; got: %v", err)
	}

	if len(code) != 6 {
		t.Fatalf("want: digit code len 6; got: %v", code)
	}

	codeInt, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		t.Fatalf("want: no error when parsing digit code to int64; got: %v", err)
	}

	if codeInt == 0 {
		t.Fatalf("want: digit code not 0; got: %v", codeInt)
	}
}
