package settings

import (
	"context"

	tablecreator "github.com/esmailemami/eshop/app/services/settings/table_creator"
	dbpkg "github.com/esmailemami/eshop/db"
)

func Initialize() {
	bindSystemSettings()
}

func bindSystemSettings() {
	baseDB := dbpkg.MustGormDBConn(context.Background())
	db, _ := baseDB.DB()
	if err := tablecreator.CreateOrUpdate(db, &SystemSetting{}); err != nil {
		panic(err)
	}
	systemSettingBinded, err := tablecreator.Bind[SystemSetting](db, SystemSetting{})
	if err != nil {
		panic(err)
	}
	systemSetting = systemSettingBinded
}
