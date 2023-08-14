package settings

import (
	"context"
	"fmt"

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
	err := tablecreator.Bind[SystemSetting](db, systemSetting)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", systemSetting)
}
