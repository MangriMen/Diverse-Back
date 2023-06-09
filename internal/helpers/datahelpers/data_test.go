// Package datahelpers provides functions to work with media or document files.
package datahelpers

import (
	"mime/multipart"
	"testing"
)

func TestValidateContentType(t *testing.T) {
	type args struct {
		fileheader           *multipart.FileHeader
		allowedMIMEBaseTypes []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateContentType(tt.args.fileheader, tt.args.allowedMIMEBaseTypes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateContentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateContentType() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateContentType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
