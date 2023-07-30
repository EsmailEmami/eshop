package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	_, err = os.Stat(filePath)
	if err != nil {
		return
	}

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

func DeleteFileByPath(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func ValidateItem(db *gorm.DB, itemID uuid.UUID, fileType models.FileType) (multiple bool, err error) {
	switch fileType {
	case models.FileTypeSystematic:
		return true, nil
	case models.FileTypeProduct:
		if !dbpkg.Exists(db, &models.Product{}, "id = ?", itemID) {
			return true, fmt.Errorf("no item found with Id #%s", itemID.String())
		}

		return true, nil
	case models.FileTypeBrand:
		if !dbpkg.Exists(db, &models.Brand{}, "id = ?", itemID) {
			return false, fmt.Errorf("no item found with Id #%s", itemID.String())
		}

		return false, nil
	default:
		return true, fmt.Errorf("invalid file type")
	}
}

func InsertItemFile(db, tx *gorm.DB, itemID uuid.UUID, fileType models.FileType, files ...*models.File) error {
	switch fileType {
	case models.FileTypeSystematic:
		return nil
	case models.FileTypeProduct:
		mapItems := []models.ProductFileMap{}

		for _, file := range files {
			mapItems = append(mapItems, models.ProductFileMap{
				ProductID: itemID,
				FileID:    *file.ID,
			})
		}

		return tx.CreateInBatches(mapItems, len(mapItems)).Error
	case models.FileTypeBrand:

		return tx.Model(&models.Brand{}).Where("id = ?", itemID).UpdateColumn("file_id", *files[0].ID).Error
	default:
		return nil
	}
}

func DeleteFile(tx *gorm.DB, file *models.File) error {
	// the default file should not delete default file
	if file.ID.String() == consts.FILE_DEFAULT_ID {
		return nil
	}

	path := GetFilePath(file)

	if err := tx.Delete(&file).Error; err != nil {
		return err
	}

	return DeleteFileByPath(path)
}

func GetFilePath(file *models.File) string {
	return GetPath(file.FileType.GetDirectory(), file.UniqueFileName)
}

func GenrateFileWhereClause(db *gorm.DB, itemID uuid.UUID, fileType models.FileType) (*gorm.DB, bool) {
	switch fileType {
	case models.FileTypeSystematic:
		return nil, false
	case models.FileTypeProduct:
		return db.Model(&models.ProductFileMap{}).Where("product_id = ?", itemID).Select("file_id"), true
	case models.FileTypeBrand:
		return db.Model(&models.Brand{}).Where("id = ?", itemID).Select("file_id"), true
	default:
		return nil, false
	}
}
