package application

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	gameErrors "github.com/v-escobar/game-api-go/internal/application/errors"
	"github.com/v-escobar/game-api-go/internal/domain/game"
	gameMocks "github.com/v-escobar/game-api-go/internal/domain/game/mocks"
	"gorm.io/gorm"
	"testing"
)

type GameServiceTestSuite struct {
	suite.Suite
	repo        *gameMocks.MockRepository
	gameService GameService
}

func (suite *GameServiceTestSuite) SetupTest() {
	suite.repo = &gameMocks.MockRepository{}
	suite.gameService = NewGameService(suite.repo)
}

func (suite *GameServiceTestSuite) TestFindAllWithNoResults() {
	suite.repo.EXPECT().FindAll().Return([]game.Game{}, nil).Once()

	games, err := suite.gameService.FindAll()

	suite.NoError(err)
	suite.Len(games, 0)
}

func (suite *GameServiceTestSuite) TestFindAllWithResults() {
	games := []game.Game{
		{ID: 1, Title: "Game 1"},
		{ID: 2, Title: "Game 2"},
	}

	suite.repo.EXPECT().FindAll().Return(games, nil).Once()

	result, err := suite.gameService.FindAll()

	suite.NoError(err)
	suite.Len(result, 2)
	suite.Equal(games, result)
}

func (suite *GameServiceTestSuite) TestFindAllWithErrors() {
	suite.repo.EXPECT().FindAll().Return(nil, errors.Errorf("Some error")).Once()

	result, err := suite.gameService.FindAll()

	suite.Error(err)
	suite.Nil(result)
	suite.Equal("Some error", err.Error())
}

func (suite *GameServiceTestSuite) TestFindByIdWithGameNotFound() {
	suite.repo.EXPECT().FindById(uint64(1)).Return(nil, gorm.ErrRecordNotFound).Once()

	result, err := suite.gameService.FindById(1)

	suite.Error(err)
	suite.Equal(err, gameErrors.GameNotFoundError)
	suite.Nil(result)
}

func TestGameService(t *testing.T) {
	suite.Run(t, new(GameServiceTestSuite))
}
