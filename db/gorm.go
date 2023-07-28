package db

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var gormDBConn *gorm.DB

func MustGormDBConn(ctx context.Context) *gorm.DB {
	db, err := GormDBConn(ctx)
	if err != nil {
		panic(err)
	}
	return db
}

func GormDBConn(ctx context.Context) (*gorm.DB, error) {
	if gormDBConn != nil {
		return gormDBConn.WithContext(ctx), nil
	}
	user := viper.GetString("appdb.username")
	password := viper.GetString("appdb.password")
	host := viper.GetString("appdb.host")
	dbName := viper.GetString("appdb.db_name")
	port := viper.GetString("appdb.port")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable  TimeZone=Asia/Tehran",
		host, port, user, dbName, password,
	)

	dbConn, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		FullSaveAssociations: false,
		Logger:               logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	dbConn = dbConn.Omit(clause.Associations)

	loadCallbacks(dbConn)

	sqlDB, _ := dbConn.DB()
	sqlDB.SetMaxIdleConns(30)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	gormDBConn = dbConn

	return gormDBConn.WithContext(ctx), nil
}

var gormLogDBConn *gorm.DB

func GormLogDBConn(ctx context.Context) (*gorm.DB, error) {
	if gormLogDBConn != nil {
		return gormLogDBConn.WithContext(ctx), nil
	}
	user := viper.GetString("log-db.username")
	password := viper.GetString("log-db.password")
	host := viper.GetString("log-db.host")
	dbName := viper.GetString("log-db.db_name")
	port := viper.GetString("log-db.port")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable  TimeZone=Asia/Tehran",
		host, port, user, dbName, password,
	)

	dbConn, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		FullSaveAssociations: false,
	})
	if err != nil {
		return nil, err
	}
	dbConn = dbConn.Omit(clause.Associations)

	loadCallbacks(dbConn)

	sqlDB, _ := dbConn.DB()
	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	gormLogDBConn = dbConn

	return gormLogDBConn.WithContext(ctx), nil
}
