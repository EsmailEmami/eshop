package settings

import (
	"context"

	"github.com/esmailemami/eshop/app/services/settings/processor"
	dbpkg "github.com/esmailemami/eshop/db"
)

func Initialize() {
	bindSystemSettings()
}

func bindSystemSettings() {
	baseDB := dbpkg.MustGormDBConn(context.Background())
	db, _ := baseDB.DB()
	if err := processor.CreateOrUpdate(db, &SystemSetting{}); err != nil {
		panic(err)
	}
	err := processor.Bind(db, systemSetting)
	if err != nil {
		panic(err)
	}
}
