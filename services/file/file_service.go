package file

import (
	"io"
	"mime/multipart"
	"os"
)

func UploadFile(fh multipart.FileHeader, target string, userRandomName bool) (uploadedFilePath string, err error) {
	uploadedFile, err := fh.Open()
	if err != nil {
		return "", err
	}

	defer uploadedFile.Close()

	bts, err := io.ReadAll(uploadedFile)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(target, os.ModePerm)
	if err != nil {
		return "", err
	}

	fileName := ""
	if !userRandomName {
		fileName = GenerateUniqueFilename(target, fh.Filename, 1)
	} else {
		fileName = GenerateRandomFileName(fh.Filename)
	}
	filePath := target + "/" + fileName

	err = WriteFile(filePath, bts)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func ReadFile(filePath string) (bts []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	bts, err = io.ReadAll(file)
	if err != nil {
		return
	}

	return
}

func WriteFile(filePath string, data []byte) error {
	fileToWrite, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer fileToWrite.Close()
	_, err = fileToWrite.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = fileToWrite.Write(data)
	if err != nil {
		return err
	}
	return fileToWrite.Sync()
}

func GetFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0
	}

	return info.Size()
}
