package file

import (
	"errors"
	"fmt"
	"strings"

	"github.com/esmailemami/eshop/app/consts"
	dbpkg "github.com/esmailemami/eshop/db"
	"github.com/esmailemami/eshop/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetFilePhysicallyPath(file *models.File) string {
	return GetPath(file.FileType.GetDirectory(), file.UniqueFileName)
}

func ValidateItem(db *gorm.DB, itemID uuid.UUID, fileType models.FileType) (multiple bool, err error) {
	multiple, _, table, _, _, _, _, _ := fileType.GetInfo()

	if !dbpkg.ExistsTable(db, table, "id=? AND deleted_at IS NULL", itemID) {
		err = fmt.Errorf("no item found with Id #%s", itemID.String())
	}

	return
}

func InsertItemFile(db, tx *gorm.DB, itemID uuid.UUID, fileType models.FileType, files ...*models.File) error {
	multiple, hasPriority, table, mapTable, foreignColumn, fileColumn, priorityColumn, _ := fileType.GetInfo()

	if multiple {
		var (
			mapItems     = []map[string]interface{}{}
			lastPriority = 0
		)

		if hasPriority {
			if err := db.Table(mapTable).Where(foreignColumn+"=?", itemID).Order(priorityColumn + " DESC").
				Select(priorityColumn).Limit(1).Find(&lastPriority).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
		}

		for _, file := range files {

			mapItem := map[string]interface{}{
				foreignColumn: itemID,
				fileColumn:    *file.ID,
			}

			if hasPriority {
				lastPriority++
				mapItem[priorityColumn] = lastPriority
			}

			mapItems = append(mapItems, mapItem)
		}

		return tx.Table(mapTable).CreateInBatches(mapItems, len(mapItems)).Error
	}

	return tx.Table(table).Where("id = ?", itemID).UpdateColumn(fileColumn, *files[0].ID).Error
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

	// if the file is not removeable we should not delete it physically
	if file.FileType.IsRemoveable() {
		return DeleteFileByPath(path)
	}
	return nil
}

func ChangeFilePriority(db, tx *gorm.DB, itemID, fileID uuid.UUID, fileType models.FileType, priority int) error {
	if priority < 0 {
		return errors.New(consts.InvalidPriority)
	}

	multiple, hasPriority, _, mapTable, foreignColumn, fileColumn, priorityColumn, _ := fileType.GetInfo()

	if !multiple || !hasPriority {
		return nil
	}

	// ignore it
	if dbpkg.ExistsTable(db, mapTable, generateStrWhere(priorityColumn, foreignColumn, fileColumn), priority, itemID, fileID) {
		return nil
	}

	// move +1 existed items priority
	if dbpkg.ExistsTable(db, mapTable, generateStrWhere(priorityColumn, foreignColumn), priority, itemID) {
		if err := tx.Table(mapTable).
			Where(foreignColumn+"=? AND "+priorityColumn+">=?", itemID, priority).
			Update(priorityColumn, gorm.Expr(priorityColumn+" + 1")).Error; err != nil {
			return err
		}
	}

	// update selected file priority
	if err := tx.Table(mapTable).
		Where(generateStrWhere(foreignColumn, fileColumn), itemID, fileID).
		Update(priorityColumn, priority).Error; err != nil {
		return err
	}

	return nil
}

func generateStrWhere(columns ...string) string {
	placeholders := make([]string, len(columns))
	for i, col := range columns {
		placeholders[i] = col + "=?"
	}
	return strings.Join(placeholders, " AND ")
}
