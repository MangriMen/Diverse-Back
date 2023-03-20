package configs

const BodyLimit = 1024 * 1024 * 1024

const PasswordEncryptCost = 12

const DataPath = "/var/lib/backend-data/"

const (
	MIMEBaseImage = "image"
)

var UploadMIMEBaseTypes = []string{
	MIMEBaseImage,
}

const WebpQuality = 75
