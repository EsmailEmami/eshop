package file

import (
	"encoding/xml"
	"mime/multipart"
)

func IsVoiceFileMimeType(mimeType string) bool {
	switch mimeType {
	case "audio/mpeg", "audio/flac", "audio/mp4", "audio/ogg", "audio/wav", "audio/aac":
		return true
	default:
		return false
	}
}

func IsPdfFileMimeType(mimeType string) bool {
	switch mimeType {
	case "application/pdf":
		return true
	default:
		return false
	}
}

func IsExcelFileMimeType(mimeType string) bool {
	switch mimeType {
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return true
	default:
		return false
	}
}

func IsSvgFileMimeType(mimeType string) bool {
	switch mimeType {
	case "image/svg+xml":
		return true
	default:
		return false
	}
}

// check given content that it can be parsed as a xml/svg
func IsValidUploadedSvgFile(fileToUpload multipart.FileHeader) bool {
	ct := fileToUpload.Header.Get("Content-Type")
	if !IsSvgFileMimeType(ct) {
		return false
	}

	uploadedFile, err := fileToUpload.Open()
	if err != nil {
		return false
	}

	defer uploadedFile.Close()
	var svg any
	err = xml.NewDecoder(uploadedFile).Decode(&svg)
	if err != nil {
		return false
	}
	return true
}
