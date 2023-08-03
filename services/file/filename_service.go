package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

func GetFileExetension(filename string) string {
	p := strings.Split(filename, ".")
	if len(p) <= 1 {
		return ""
	}

	return strings.ToLower(p[len(p)-1])
}

// GenerateUniqueFilename creates a unique name for the requested file in the specified directory to prevent name clashes and overwriting of old files.
func GenerateUniqueFilename(dir, orginalFilename string, duplicationIndex uint) string {
	if duplicationIndex == 1 {
		filePath := fmt.Sprintf("%s/%s", dir, orginalFilename)
		if _, err := os.Stat(filePath); err != nil {
			// file not exists
			return orginalFilename
		}
	}
	ext := GetFileExetension(orginalFilename)
	newFilename := strings.TrimRight(orginalFilename, "."+ext)

	filePath := fmt.Sprintf("%s/%s-%03d.%s", dir, newFilename, duplicationIndex, ext)
	if _, err := os.Stat(filePath); err != nil {
		// file not exists
		return fmt.Sprintf("%s-%03d.%s", newFilename, duplicationIndex, ext)
	}

	duplicationIndex++
	return GenerateUniqueFilename(dir, orginalFilename, duplicationIndex)
}

func GenerateRandomFileName(orginalFilename string) string {
	ext := GetFileExetension(orginalFilename)
	return uuid.NewString() + "." + ext
}
