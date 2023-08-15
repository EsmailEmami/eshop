package settings

import (
	"context"

	"github.com/esmailemami/eshop/app/services/settings/processor"
	dbpkg "github.com/esmailemami/eshop/db"
)

var systemSetting = new(SystemSetting)

type SystemSetting struct {
	FileExpireTimeStampts *int `column:"file_expire_time_stampts" title:"File Expire TimeStampts" description:"The value is calculating by month"`
}

func (SystemSetting) TableName() string {
	return "system_setting"
}

func (SystemSetting) SchemaName() string {
	return "public"
}

func (s *SystemSetting) LoadDefaultValues() {
	s.FileExpireTimeStampts = func() *int {
		value := 4
		return &value
	}()
}

func GetSystemSettings() (*SystemSetting, error) {
	if systemSetting == nil {
		bindSystemSettings()
	}

	return systemSetting, nil
}

func UpdateSystemSettings(field string, value interface{}) error {
	baseDB := dbpkg.MustGormDBConn(context.Background())
	db, _ := baseDB.DB()
	return processor.UpdateField(db, systemSetting, field, value)
}

func GetSystemSettingItems() []processor.SettingItem {
	return processor.GetItems(systemSetting)
}
