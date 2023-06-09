package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MangriMen/Diverse-Back/configs"
	"github.com/joho/godotenv"
)

// OpenTempFile opens file in project temp directory.
func OpenTempFile(filename string) *os.File {
	file, err := os.OpenFile(
		filepath.Join(FindRoot(), configs.TempDirectory, filename),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600,
	)
	if err != nil {
		panic(err)
	}

	return file
}

// FindRoot returns the absolute path of the root directory.
// It searches for the 'go.mod' file from the current working directory upwards
func FindRoot() string {
	if IsRunningInContainer() {
		return ""
	}

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err = os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return currentDir
}

// LoadEnvironment loads the environment variables from the .env file.
func LoadEnvironment(envFile string) {
	err := godotenv.Load(dir(envFile))
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %w", err))
	}
}

// dir returns the absolute path of the given environment file (envFile) in the Go module's
// root directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the envFile to the directory containing 'go.mod'.
// It panics if it fails to find the 'go.mod' file.
func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err = os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}

// IsRunningInContainer return the status of program instance,
// running in docker container or not.
func IsRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return false
	}
	return true
}
