package application

import (
	"errors"
	applicationErrors "github.com/v-escobar/game-api-go/internal/application/errors"
	"github.com/v-escobar/game-api-go/internal/domain/game"
	"github.com/v-escobar/game-api-go/internal/interfaces/dto"
	"gorm.io/gorm"
	"time"
)

type GameService interface {
	FindAll() ([]game.Game, error)
	FindById(id uint64) (*game.Game, error)
	Create(d *dto.Game) error
}

type GameServiceImplementation struct {
	repo game.Repository
}

func NewGameService(repo game.Repository) GameService {
	return &GameServiceImplementation{repo: repo}
}

func (g *GameServiceImplementation) FindAll() ([]game.Game, error) {
	games, err := g.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (g *GameServiceImplementation) FindById(id uint64) (*game.Game, error) {
	persistedGame, err := g.repo.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = applicationErrors.GameNotFoundError
		}
		return nil, err
	}
	return persistedGame, nil
}

func (g *GameServiceImplementation) Create(d *dto.Game) error {
	gameToPersist := game.Game{
		Title:     d.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := g.repo.Create(&gameToPersist); err != nil {
		return err
	}

	d.ID = gameToPersist.ID
	return nil
}
