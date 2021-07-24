package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ewenquim/horkruxes/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbOptions struct {
	test bool
}

func initDatabase(options dbOptions) *gorm.DB {
	db_name := "db.sqlite3"
	if options.test {
		db_name = "db_test.sqlite3"
		os.Remove(db_name)
	}
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	// Opens Database
	db, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
	db.AutoMigrate(&model.Message{})
	fmt.Println("Database Migrated")
	return db
}
