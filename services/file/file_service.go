package file

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadFile(fh *multipart.FileHeader, target string, useRandomName, startFromRoot bool) (uploadedFilePath, fileName string, err error) {
	uploadedFile, err := fh.Open()
	if err != nil {
		return "", "", err
	}

	defer uploadedFile.Close()

	bts, err := io.ReadAll(uploadedFile)
	if err != nil {
		return "", "", err
	}

	var dirPath string

	if startFromRoot {
		dirPath = GetPath(target)
	} else {
		dirPath = target
	}

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return "", "", err
	}

	if !useRandomName {
		fileName = GenerateUniqueFilename(target, fh.Filename, 1)
	} else {
		fileName = GenerateRandomFileName(fh.Filename)
	}

	filePath := dirPath + "/" + fileName

	err = WriteFile(filePath, bts)
	if err != nil {
		return "", "", err
	}

	return filePath, fileName, nil
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

// GetPath joins the given path with root
func GetPath(paths ...string) string {
	// Get the current working directory (project's root directory)
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err) // Handle the error appropriately based on your use case
	}

	// Concatenate the current directory with "/upload" to get the upload folder path
	paths = append([]string{currentDir}, paths...)
	pathDir := filepath.Join(paths...)

	return pathDir
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
