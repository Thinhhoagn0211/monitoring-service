package serializer

import (
	"strings"

	"github.com/chchench/textract"
)

func GetContentFile(path string) (string, error) {
	textract, err := textract.RetrieveTextFromFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(textract), nil
}
