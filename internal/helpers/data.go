package helpers

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"golang.org/x/exp/slices"
)

func GenerateUniqueFilename() string {
	uniqueFileId := uuid.New()
	filename := strings.Replace(uniqueFileId.String(), "-", "", -1)
	return filename
}

func ParseContentType(contentType string) (string, string) {
	parsedContentType := strings.Split(contentType, "/")
	baseType, extendType := parsedContentType[0], parsedContentType[1]
	return baseType, extendType
}

func ValidateContentType(fileheader *multipart.FileHeader, allowedMIMEBaseTypes []string) (string, string, error) {
	baseType, extendType := ParseContentType(fileheader.Header.Get("Content-Type"))
	if !slices.Contains(allowedMIMEBaseTypes, baseType) {
		return "", "", fmt.Errorf("this type of content cannot be uploaded")
	}

	return baseType, extendType, nil
}

func ProcessFile(fileheader *multipart.FileHeader, filepath string) error {
	baseType, _ := ParseContentType(fileheader.Header.Get("Content-Type"))

	file, err := fileheader.Open()
	if err != nil {
		return err
	}
	var fileBuffer []byte
	file.Read(fileBuffer)
	file.Close()

	switch baseType {
	case configs.MIMEBaseImage:
		convertedFile, err := bimg.NewImage(fileBuffer).Convert(bimg.WEBP)
		if err != nil {
			return err
		}

		processed, err := bimg.NewImage(convertedFile).Process(bimg.Options{Quality: 75})
		if err != nil {
			return err
		}

		err = bimg.Write(filepath, processed)
		if err != nil {
			return err
		}
	}

	return nil
}
