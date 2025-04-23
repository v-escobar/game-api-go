package persistence

import (
	"github.com/v-escobar/game-api-go/internal/domain/game"
	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func (r *GameRepository) Create(g *game.Game) error {
	return r.db.Create(g).Error
}

func (r *GameRepository) FindAll() ([]game.Game, error) {
	var games []game.Game
	err := r.db.Find(&games).Error

	return games, err
}

func (r *GameRepository) FindById(id uint64) (*game.Game, error) {
	var u game.Game
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func NewGameRepository(db *gorm.DB) game.Repository {
	return &GameRepository{db}
}
