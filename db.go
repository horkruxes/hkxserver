package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/horkruxes/hkxserver/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbOptions struct {
	test bool
}

func initDatabase(db_name string) *gorm.DB {

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
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		panic("failed to migrate database")
	}
	fmt.Println("Database Migrated")
	return db
}
