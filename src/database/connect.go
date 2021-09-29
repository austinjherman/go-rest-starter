package database

import (
	userModels "aherman/src/models/user"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Init todo
func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", os.Getenv("DATABASE_NAME"))), &gorm.Config{})
	if err != nil {
  	panic("failed to connect database")
	}

	// Migrate the schema
  db.AutoMigrate(&userModels.User{})
	
	return db
}