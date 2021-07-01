package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ewenquim/horkruxes-client/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)



func initDatabase() *gorm.DB {
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
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{
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
