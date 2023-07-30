package file

import (
	"encoding/xml"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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

func IsImageMimeType(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/png", "image/gif":
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

func IsImageFile(fh *multipart.FileHeader) bool {
	contentType := fh.Header.Get("Content-Type")
	if !IsImageMimeType(contentType) {
		return false
	}

	uploadedFile, err := fh.Open()
	if err != nil {
		return false
	}
	defer uploadedFile.Close()

	_, _, err = image.Decode(uploadedFile)
	return err == nil
}

// check given content that it can be parsed as a xml/svg
func IsValidUploadedSvgFile(fileToUpload *multipart.FileHeader) bool {
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

func GetMimeType(fh *multipart.FileHeader) string {
	return fh.Header.Get("Content-Type")
}
