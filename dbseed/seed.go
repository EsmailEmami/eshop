package dbseed

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/esmailemami/eshop/db"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type seed struct {
	fn   func(conn *gorm.DB) error
	conn *gorm.DB
}

func Run() error {
	conn, err := db.GormDBConn(context.Background())
	if err != nil {
		return err
	}

	conn = conn.Session(&gorm.Session{Logger: logger.Discard})

	seeds := []seed{
		{fn: seedFile, conn: conn},
		{fn: seedRole, conn: conn},
		{fn: seedUser, conn: conn},
		{fn: seedBrand, conn: conn},
		{fn: seedColor, conn: conn},
	}

	for _, s := range seeds {
		fmt.Printf("Seeding .... %s\n", getFnName(s.fn))
		err := s.fn(s.conn)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Seeded ..... %s\n", getFnName(s.fn))
		fmt.Println("-------------------------------------------")
	}

	return nil
}

func getFnName(fn func(*gorm.DB) error) string {
	fullname := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	name := strings.Split(fullname, ".")
	return name[len(name)-1]
}
