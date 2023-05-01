// Package datahelpers provides functions to work with media or document files.
package datahelpers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"golang.org/x/exp/slices"
)

// GenerateUniqueFilename generates uuid, removes dashed from it and returns.
func GenerateUniqueFilename() string {
	uniqueFileID := uuid.New()
	filename := strings.ReplaceAll(uniqueFileID.String(), "-", "")
	return filename
}

// ParseContentType returns content type splitted by / to tuple.
func ParseContentType(contentType string) (string, string) {
	parsedContentType := strings.Split(contentType, "/")
	baseType, extendType := parsedContentType[0], parsedContentType[1]
	return baseType, extendType
}

// ValidateContentType checks that filetype in allowed types.
// Returns splitted type if ok, else error.
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

	defer helpers.CloseQuietly(file)

	fileBuffer := bytes.NewBuffer(nil)
	if _, err = io.Copy(fileBuffer, file); err != nil {
		return nil, err
	}

	return fileBuffer.Bytes(), nil
}

// ProcessFile check base file type and call process function to save it.
func ProcessFile(fileheader *multipart.FileHeader, filepath string) error {
	fileBuffer, err := readFileFromMultipart(fileheader)
	if err != nil {
		return err
	}

	baseType, _ := ParseContentType(fileheader.Header.Get("Content-Type"))

	if baseType == configs.MIMEBaseImage {
		if err = processImage(fileBuffer, configs.WebpQuality, filepath); err != nil {
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

	if err = bimg.Write(filepath, processed); err != nil {
		return err
	}

	return nil
}
