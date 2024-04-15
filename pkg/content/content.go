package content

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func Base64Type(contentEncoded string) ([]byte, string, error) {
	fileContentDecoded, err := base64.StdEncoding.DecodeString(contentEncoded)
	if err != nil {
		return nil, "", err
	}
	contentType := http.DetectContentType(fileContentDecoded)
	return fileContentDecoded, strings.Split(contentType, ";")[0], nil
}
