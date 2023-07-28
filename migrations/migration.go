package migrations

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/esmailemami/eshop/db"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed *yaml
var migrationFS embed.FS

var dbConn *gorm.DB
var logdbConn *gorm.DB

func getDB() *gorm.DB {
	if dbConn != nil {
		return dbConn
	}
	var err error
	dbConn, err = db.GormDBConn(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	dbConn = dbConn.Session(&gorm.Session{SkipHooks: false, NewDB: false, Logger: logger.Discard})
	return dbConn
}

func getLogDB() *gorm.DB {
	if logdbConn != nil {
		return logdbConn
	}
	conn, err := db.GormLogDBConn(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	logdbConn = conn.Session(&gorm.Session{SkipHooks: false, NewDB: false, Logger: logger.Discard})

	return logdbConn
}

func init() {
	log.SetFlags(log.Lshortfile)
}

type migrationYamlFile struct {
	Up         string `yaml:"up"`
	Down       string `yaml:"down"`
	Connection string `yaml:"connection"`
}

type migration struct {
	ID       uint
	Name     string
	Batch    uint
	Filename string
}

func (m migration) Down() error {
	dbConn := getDB()
	logdbConn := getLogDB()

	name := m.Name
	fmt.Println("Rollback ", name)
	filename := m.Filename

	bts, err := migrationFS.ReadFile(filename)
	if err != nil {
		return err
	}

	var mf migrationYamlFile
	err = yaml.Unmarshal(bts, &mf)
	if err != nil {
		return err
	}
	migrateSql := mf.Down

	var tx *gorm.DB

	switch mf.Connection {
	case "log-db":
		tx = logdbConn.Begin()
	default:
		tx = dbConn.Begin()
	}

	err = tx.Exec(migrateSql).Error
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	err = dbConn.Delete(&m).Error
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

type migrations struct{}

// MakeMigration create new migration file
func MakeMigration(create *string) error {
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("./migrations/%v_%v.yaml", timestamp, strcase.ToSnake(*create))
	content := `---
up: |
  -- UP SQL

down: |
  -- DOWN SQL
`
	f, e := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
	if e != nil {
		fmt.Println(e)
		return e
	}
	defer f.Close()
	_, err := f.WriteString(content)
	return err
}

// checkMigrationTable will check existence of migrations table and create it if it doesn't exist.
func checkMigrationTable() {
	dbConn := getDB()

	dbConn.Exec(`CREATE SCHEMA IF NOT EXISTS public;`)
	dbConn.Exec(`CREATE TABLE IF NOT EXISTS migrations (
    id serial NOT NULL,
    name varchar(255) NOT NULL,
    batch integer,
    filename varchar(512) NOT NULL,
    CONSTRAINT migrations_pkey PRIMARY KEY (id)
    )`)
}

// Migrate call Migrate function of files and save batches in database
func Migrate() error {
	checkMigrationTable()

	dbConn := getDB()
	logdbConn := getLogDB()

	files, e := migrationFS.ReadDir(".")
	if e != nil {
		return errors.New("Error on loading migrations directory")
	}

	var batch struct {
		LastBatch uint `sql:"last_batch"`
	}
	dbConn.Raw(`select max(batch) as last_batch from migrations;`).Scan(&batch)
	batch.LastBatch++

	for _, v := range files {
		if v.IsDir() {
			continue
		}
		var mg migration

		filename := v.Name()
		migrationName := strings.TrimSuffix(filename, ".yml")
		migrationName = strings.TrimSuffix(migrationName, ".yaml")
		underscoreIndex := strings.Index(filename, "_")
		migrationName = migrationName[underscoreIndex+1:]
		fmt.Println("Migrating", migrationName)

		_ = dbConn.Where("name LIKE ?", migrationName).First(&mg).Error

		if mg.ID != 0 {
			//This file already migrated
			continue
		}
		file, err := migrationFS.Open(v.Name())
		if err != nil {
			return err
		}

		var mf migrationYamlFile
		err = yaml.NewDecoder(file).Decode(&mf)
		if err != nil {
			return err
		}

		var tx *gorm.DB

		migrateSql := mf.Up
		switch mf.Connection {
		case "log-db":
			tx = logdbConn.Begin()
		default:
			tx = dbConn.Begin()
		}

		err = tx.Exec(migrateSql).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}

		// save migrated file to DB
		mg.Name = migrationName
		mg.Batch = batch.LastBatch
		mg.Filename = filename
		batch.LastBatch = batch.LastBatch + 1
		err = dbConn.Create(&mg).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}

		tx.Commit()

		fmt.Println("Migrated ", migrationName)
	}

	return dbConn.Exec("update migrations set batch=coalesce((select max(batch) from migrations) , 0)+1 where batch is null;").Error
}

// Rollback will rollback database using batch number
func Rollback() error {
	checkMigrationTable()
	dbConn := getDB()

	var rollbacks []migration
	dbConn.Where("batch = (select max(batch) from migrations)").Find(&rollbacks)
	for _, v := range rollbacks {
		err := v.Down()
		if err != nil {
			return err
		}
	}

	return nil
}

func RollbackAll() error {
	checkMigrationTable()
	dbConn := getDB()

	var rollbacks []migration
	dbConn.Order("id desc").Find(&rollbacks)
	for _, v := range rollbacks {
		err := v.Down()
		if err != nil {
			return err
		}
	}

	return nil
}
