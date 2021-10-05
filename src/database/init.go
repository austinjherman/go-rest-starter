package database

import (
	tokenModels "aherman/src/models/token"
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
  db.AutoMigrate(&tokenModels.Token{})
  db.AutoMigrate(&userModels.User{})
	
	return db
}