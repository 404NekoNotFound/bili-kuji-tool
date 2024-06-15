package db

import (
	"bili-kuji-management/src/db/table"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DB struct {
	db  *gorm.DB
	log *zap.Logger
}

// Connect to the database
func Connect(log *zap.Logger) *DB {
	db, err := gorm.Open(sqlite.Open("kuji.db?_pragma=journal_mode(WAL)"), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.AutoMigrate(&table.Account{}, &table.Reward{}, &table.Price{}); err != nil {
		log.Fatal(err.Error())
	}

	return &DB{db, log}
}
