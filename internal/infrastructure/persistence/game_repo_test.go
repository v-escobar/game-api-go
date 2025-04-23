package persistence

import (
	"errors"
	"github.com/v-escobar/game-api-go/internal/infrastructure/db"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/v-escobar/game-api-go/internal/domain/game"
	"gorm.io/gorm"
)

type GameRepositoryTestSuite struct {
	suite.Suite
	repo game.Repository
	db   *gorm.DB
	game *game.Game
}

func (suite *GameRepositoryTestSuite) SetupSuite() {
	postgresContainer, _ := postgres.Run(suite.T().Context(),
		"postgres:16",
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	dsn, _ := postgresContainer.ConnectionString(suite.T().Context())

	dbInstance := db.NewDB(dsn)

	suite.db = dbInstance
	suite.repo = NewGameRepository(dbInstance)
}

func (suite *GameRepositoryTestSuite) BeforeTest(string, string) {
	suite.game = &game.Game{Title: "Test Game", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	suite.db.Create(suite.game)
}

func (suite *GameRepositoryTestSuite) AfterTest(string, string) {
	suite.db.Exec("DELETE FROM games")
}

func (suite *GameRepositoryTestSuite) TestLoadAllGames() {
	games, _ := suite.repo.FindAll()

	suite.Len(games, 1)
	suite.Equal(games[0].ID, suite.game.ID)
	suite.Equal(games[0].Title, suite.game.Title)
}

func (suite *GameRepositoryTestSuite) TestLoadGameById() {
	persistedGame, _ := suite.repo.FindById(uint64(suite.game.ID))

	suite.EqualExportedValues(persistedGame, suite.game)
}

func (suite *GameRepositoryTestSuite) TestLoadGameByIdNotFound() {
	persistedGame, err := suite.repo.FindById(999999)

	suite.Nil(persistedGame)
	suite.Error(err)
	suite.Equal(errors.Is(err, gorm.ErrRecordNotFound), true)
}

func (suite *GameRepositoryTestSuite) TestCreateGame() {
	var newGame = game.Game{Title: "New Game", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err := suite.repo.Create(&newGame)

	suite.NoError(err)
	suite.NotZero(newGame.ID)
}

func TestGamesRepository(t *testing.T) {
	suite.Run(t, new(GameRepositoryTestSuite))
}
