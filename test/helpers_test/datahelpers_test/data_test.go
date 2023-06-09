// Package datahelpers_test provides test for datahelpers package.
package datahelpers_test

import (
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/MangriMen/Diverse-Back/internal/helpers/datahelpers"
)

func TestParseContentType(t *testing.T) {
	tests := []struct {
		name               string
		contentType        string
		expectedBaseType   string
		expectedExtendType string
	}{
		{
			name:               "Test divide types of image/png",
			contentType:        "image/png",
			expectedBaseType:   "image",
			expectedExtendType: "png",
		},
		{
			name:               "Test divide types of video/mp4",
			contentType:        "video/mp4",
			expectedBaseType:   "video",
			expectedExtendType: "mp4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := datahelpers.ParseContentType(tt.contentType)
			if got != tt.expectedBaseType {
				t.Errorf("ParseContentType() got = %v, want %v", got, tt.expectedBaseType)
			}
			if got1 != tt.expectedExtendType {
				t.Errorf("ParseContentType() got1 = %v, want %v", got1, tt.expectedExtendType)
			}
		})
	}
}

func TestValidateContentType(t *testing.T) {
	type args struct {
		fileheader           *multipart.FileHeader
		allowedMIMEBaseTypes []string
	}
	tests := []struct {
		name               string
		args               args
		expectedBaseType   string
		expectedExtendType string
		wantErr            bool
	}{
		{
			name: "Test image/png is allowed",
			args: args{
				fileheader:           &multipart.FileHeader{Header: textproto.MIMEHeader{"Content-Type": []string{"image/png"}}},
				allowedMIMEBaseTypes: []string{"image"},
			},
			expectedBaseType:   "image",
			expectedExtendType: "png",
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := datahelpers.ValidateContentType(tt.args.fileheader, tt.args.allowedMIMEBaseTypes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateContentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expectedBaseType {
				t.Errorf("ValidateContentType() got = %v, want %v", got, tt.expectedBaseType)
			}
			if got1 != tt.expectedExtendType {
				t.Errorf("ValidateContentType() got1 = %v, want %v", got1, tt.expectedExtendType)
			}
		})
	}
}
