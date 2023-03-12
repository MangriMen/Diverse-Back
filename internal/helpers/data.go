package helpers

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func GenerateUniqueFilename() string {
	uniqueFileId := uuid.New()
	filename := strings.Replace(uniqueFileId.String(), "-", "", -1)
	return filename
}

func ValidateContentType(fileheader *multipart.FileHeader, allowedMIMEBaseTypes []string) (string, string, error) {
	contentType := strings.Split(fileheader.Header.Get("Content-Type"), "/")

	baseType, extendType := contentType[0], contentType[1]
	if !slices.Contains(allowedMIMEBaseTypes, baseType) {
		return "", "", fmt.Errorf("this type of content cannot be uploaded")
	}

	return baseType, extendType, nil
}
