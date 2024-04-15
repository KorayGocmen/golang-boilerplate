package content

import (
	"encoding/base64"
	"testing"
)

func TestBase64Type(t *testing.T) {
	tt := []struct {
		content     string
		contentType string
	}{
		{
			// 1px image.
			content:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII=",
			contentType: "image/png",
		},
		{
			// text.
			content:     base64.StdEncoding.EncodeToString([]byte("hello world")),
			contentType: "text/plain",
		},
		{
			// smallest pdf.
			content:     "JVBERi0xLjIgCjkgMCBvYmoKPDwKPj4Kc3RyZWFtCkJULyAzMiBUZiggIFlPVVIgVEVYVCBIRVJFICAgKScgRVQKZW5kc3RyZWFtCmVuZG9iago0IDAgb2JqCjw8Ci9UeXBlIC9QYWdlCi9QYXJlbnQgNSAwIFIKL0NvbnRlbnRzIDkgMCBSCj4+CmVuZG9iago1IDAgb2JqCjw8Ci9LaWRzIFs0IDAgUiBdCi9Db3VudCAxCi9UeXBlIC9QYWdlcwovTWVkaWFCb3ggWyAwIDAgMjUwIDUwIF0KPj4KZW5kb2JqCjMgMCBvYmoKPDwKL1BhZ2VzIDUgMCBSCi9UeXBlIC9DYXRhbG9nCj4+CmVuZG9iagp0cmFpbGVyCjw8Ci9Sb290IDMgMCBSCj4+CiUlRU9G",
			contentType: "application/pdf",
		},
	}

	for _, tc := range tt {
		_, got, err := Base64Type(tc.content)
		if err != nil {
			t.Errorf("want: err nil; got: err = %v", err)
		}

		if got != tc.contentType {
			t.Errorf("want: content type %v; got: %v", tc.contentType, got)
		}
	}
}
