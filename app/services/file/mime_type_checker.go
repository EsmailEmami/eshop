package file

import (
	"encoding/xml"
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
	case "image/jpeg", "image/png", "image/gif", "application/octet-stream":
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

	// _, _, err = image.Decode(uploadedFile)
	// return err == nil

	buffer := make([]byte, 512)
	n, err := uploadedFile.Read(buffer)
	if err != nil {
		return false
	}

	// Check if the data appears to be an image
	if n >= 3 && buffer[0] == 0xFF && buffer[1] == 0xD8 && buffer[2] == 0xFF {
		// JPEG magic number
		return true
	} else if n >= 4 && buffer[0] == 0x89 && buffer[1] == 0x50 && buffer[2] == 0x4E && buffer[3] == 0x47 {
		// PNG magic number
		return true
	} else if n >= 6 && buffer[0] == 0x47 && buffer[1] == 0x49 && buffer[2] == 0x46 && buffer[3] == 0x38 && (buffer[4] == 0x37 || buffer[4] == 0x39) && buffer[5] == 0x61 {
		// GIF87a or GIF89a magic number
		return true
	} else if n >= 12 && buffer[8] == 'W' && buffer[9] == 'E' && buffer[10] == 'B' && buffer[11] == 'P' {
		// WebP magic number
		return true
	}

	return false
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
	return err == nil
}

func GetMimeType(fh *multipart.FileHeader) string {
	return fh.Header.Get("Content-Type")
}
