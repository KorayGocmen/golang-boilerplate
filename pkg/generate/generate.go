package generate

import (
	cryptorand "crypto/rand"
	"fmt"
	mathrand "math/rand"
	"strings"
	"time"
)

func AlphaCode(length int, includeNumbers bool) string {
	random := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	alphabet := "ABCDEFGHIJKLMNPQRSTUVWXYZ" // O is omitted not to be confused with 0.
	if includeNumbers {
		alphabet += "123456789"
	}

	var chars []byte
	for i := 0; i < length; i++ {
		chars = append(chars, alphabet[random.Intn(len(alphabet))])
	}

	return string(chars)
}

func DigitCode(length int) (string, error) {
	digits := make([]byte, length)
	if _, err := cryptorand.Read(digits); err != nil {
		err = fmt.Errorf("rand read error: %w", err)
		return "", err
	}

	for i := 0; i < length; i++ {
		digits[i] = uint8(48 + (digits[i] % 10))
	}

	return string(digits), nil
}

// Filename returns a filename safe string.
func FileName(prefix, ext string, length int, includeTs bool) string {
	if includeTs {
		now := time.Now().UTC().Format(time.RFC3339)
		prefix += strings.ReplaceAll(now, ":", "")
		prefix += "-"
	}

	prefix = strings.TrimSpace(strings.ToLower(prefix))
	name := strings.ToLower(AlphaCode(length, false))
	ext = strings.TrimSpace(strings.ToLower(ext))

	return fmt.Sprintf("%s%s%s", prefix, name, ext)
}
