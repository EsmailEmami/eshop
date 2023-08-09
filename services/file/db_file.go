package file

import (
	"errors"
	"fmt"

	"github.com/esmailemami/eshop/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetFilePhysicallyPath(file *models.File) string {
	return GetPath(file.FileType.GetDirectory(), file.UniqueFileName)
}

func ValidateItem(db *gorm.DB, itemID uuid.UUID, fileType models.FileType) (multiple bool, err error) {
	multiple, _, table, _, _, _, _ := fileType.GetInfo()

	if !dbpkg.ExistsTable(db, table, "id=? AND deleted_at IS NULL", itemID) {
		err = fmt.Errorf("no item found with Id #%s", itemID.String())
	}

	return
}

func InsertItemFile(db, tx *gorm.DB, itemID uuid.UUID, fileType models.FileType, files ...*models.File) error {
	multiple, hasPriority, table, mapTable, foreignColumn, fileColumn, _ := fileType.GetInfo()

	if multiple {
		var (
			mapItems     = []map[string]interface{}{}
			lastPriority = 0
		)

		if hasPriority {
			if err := db.Table(mapTable).Where(foreignColumn+"=? AND deleted_at IS NULL", itemID).
				Select("priority").Order("priority DESC").First(&lastPriority).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
		}

		for _, file := range files {

			mapItem := map[string]interface{}{
				foreignColumn: itemID,
				fileColumn:    *file.ID,
				"id":          uuid.New(),
			}

			if hasPriority {
				lastPriority++
				mapItem["priority"] = lastPriority
			}

			mapItems = append(mapItems, mapItem)
		}

		return tx.Table(mapTable).CreateInBatches(mapItems, len(mapItems)).Error
	}

	return tx.Model(table).Where("id = ?", itemID).UpdateColumn(fileColumn, *files[0].ID).Error
}

func GenrateFileWhereClause(
	db *gorm.DB,
	itemID uuid.UUID,
	fileType models.FileType,
) (*gorm.DB, bool) {
	switch fileType {
	case models.FileTypeSystematic:
		return nil, false
	case models.FileTypeProduct:
		return db.Model(&models.ProductFileMap{}).
				Where("product_id = ?", itemID).
				Order("priority ASC").
				Select("file_id"),
			true
	case models.FileTypeBrand:
		return db.Model(&models.Brand{}).Where("id = ?", itemID).Select("file_id"), true
	case models.FileTypeAppPic:
		return db.Model(&models.AppPic{}).Where("id = ?", itemID).Select("file_id"), true
	default:
		return nil, false
	}
}

func DeleteFile(tx *gorm.DB, file *models.File) error {
	// the default file should not delete default file
	if file.ID.String() == consts.FILE_DEFAULT_ID {
		return nil
	}

	path := GetFilePhysicallyPath(file)

	if err := tx.Delete(file).Error; err != nil {
		return err
	}

	return DeleteFileByPath(path)
}

func ChangeFilePriority(db, tx *gorm.DB, itemID, fileID uuid.UUID, fileType models.FileType, priority int) error {
	if priority < 0 {
		return errors.New(consts.InvalidPriority)
	}

	switch fileType {
	case models.FileTypeProduct:

		if dbpkg.Exists(db, &models.ProductFileMap{}, "priority=? AND product_id=? AND file_id=?", priority, itemID, fileID) {
			return nil
		}

		// update existed items priority
		if dbpkg.Exists(db, &models.ProductFileMap{}, "priority=? AND product_id=?", priority, itemID) {
			if err := tx.Model(&models.ProductFileMap{}).
				Where("product_id = ? AND priority >= ?", itemID, priority).
				Update("priority", gorm.Expr("priority + 1")).Error; err != nil {
				return err
			}
		}

		// update selected file priority
		if err := tx.Model(&models.ProductFileMap{}).
			Where("product_id = ? AND file_id >= ?", itemID, fileID).
			Update("priority", priority).Error; err != nil {
			return err
		}

		return nil
	default:
		return nil
	}
}
