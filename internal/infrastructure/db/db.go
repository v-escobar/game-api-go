package db

import (
	"github.com/v-escobar/game-api-go/internal/domain/game"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewDB(dsn string) *gorm.DB {
	dbInstance, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Fatal(err)
	}

	_ = dbInstance.AutoMigrate(&game.Game{})

	return dbInstance
}
