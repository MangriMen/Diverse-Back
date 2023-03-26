package helpers

import (
	"bytes"
	"fmt"
	"io"
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

func ValidateContentType(
	fileheader *multipart.FileHeader,
	allowedMIMEBaseTypes []string,
) (string, string, error) {
	baseType, extendType := ParseContentType(fileheader.Header.Get("Content-Type"))
	if !slices.Contains(allowedMIMEBaseTypes, baseType) {
		return "", "", fmt.Errorf("this type of content cannot be uploaded")
	}

	return baseType, extendType, nil
}

func readFileFromMultipart(fileheader *multipart.FileHeader) ([]byte, error) {
	file, err := fileheader.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBuffer, file); err != nil {
		return nil, err
	}

	return fileBuffer.Bytes(), nil
}

func ProcessFile(fileheader *multipart.FileHeader, filepath string) error {
	fileBuffer, err := readFileFromMultipart(fileheader)
	if err != nil {
		return err
	}

	baseType, _ := ParseContentType(fileheader.Header.Get("Content-Type"))
	switch baseType {
	case configs.MIMEBaseImage:
		if err := processImage(fileBuffer, 85, filepath); err != nil {
			return err
		}
	}

	return nil
}

func processImage(buffer []byte, quality int, filepath string) error {
	convertedFile, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		return err
	}

	processed, err := bimg.NewImage(convertedFile).Process(bimg.Options{Quality: quality})
	if err != nil {
		return err
	}

	err = bimg.Write(filepath, processed)
	if err != nil {
		return err
	}

	return nil
}
