package game

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/v-escobar/game-api-go/internal/application"
	gameMocks "github.com/v-escobar/game-api-go/internal/domain/game/mocks"
	"github.com/v-escobar/game-api-go/internal/interfaces/dto"
	testing2 "github.com/v-escobar/game-api-go/internal/testing"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/v-escobar/game-api-go/internal/domain/game"
	"gorm.io/gorm"
)

type HandlerTestSuite struct {
	suite.Suite
	repo    *gameMocks.MockRepository
	echo    *echo.Echo
	handler *Handler
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.repo = gameMocks.NewMockRepository(suite.T())
	suite.echo = echo.New()
	suite.handler = NewGameHandler(application.NewGameService(suite.repo))
}

// ListGames

func (suite *HandlerTestSuite) TestListGamesEmptyResult() {
	suite.repo.EXPECT().FindAll().Return([]game.Game{}, nil).Once()
	req, _ := http.NewRequest(http.MethodGet, "/games", nil)
	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	err := suite.handler.ListGames(context)
	games := testing2.UnmarshalBody[[]game.Game](recorder)

	suite.NoError(err)
	suite.Equal(http.StatusOK, recorder.Code)
	suite.Len(games, 0)
}

func (suite *HandlerTestSuite) TestListGamesWithError() {
	suite.repo.EXPECT().FindAll().Return(nil, errors.New("some error"))
	req, _ := http.NewRequest("GET", "/games", nil)
	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	err := suite.handler.ListGames(context)
	errorDto := testing2.UnmarshalBody[dto.Error](recorder)

	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, recorder.Code)
	suite.Equal("could not fetch games, Error: some error", errorDto.Message)
}

func (suite *HandlerTestSuite) TestListGamesWithResult() {
	suite.repo.EXPECT().FindAll().Return([]game.Game{{ID: 1, Title: "Test Game"}}, nil)
	req, _ := http.NewRequest(http.MethodGet, "/games", nil)
	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	err := suite.handler.ListGames(context)
	games := testing2.UnmarshalBodyArray[[]game.Game](recorder)

	suite.NoError(err)
	suite.Equal(http.StatusOK, recorder.Code)

	suite.Len(games, 1)
	suite.Equal(game.Game{ID: 1, Title: "Test Game"}, games[0])
}

// FindById

func (suite *HandlerTestSuite) TestFindByIdWithNoResult() {
	suite.repo.EXPECT().FindById(uint64(1)).Return(nil, gorm.ErrRecordNotFound)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	context.SetPath("/games/:id")
	context.SetParamNames("id")
	context.SetParamValues("1")
	err := suite.handler.GetGame(context)

	suite.NoError(err)
	suite.Equal(http.StatusNotFound, recorder.Code)
}

func (suite *HandlerTestSuite) TestFindByIdInvalidId() {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	context.SetPath("/games/:id")
	context.SetParamNames("id")
	context.SetParamValues("NaN")

	err := suite.handler.GetGame(context)
	suite.NoError(err)

	errorDto := testing2.UnmarshalBody[dto.Error](recorder)
	suite.Equal(http.StatusBadRequest, recorder.Code)
	suite.Equal("error parsing identifier: strconv.ParseUint: parsing \"NaN\": invalid syntax", errorDto.Message)
	suite.repo.AssertNotCalled(suite.T(), "FindById")
}

func (suite *HandlerTestSuite) TestFindByIdWithResult() {
	suite.repo.EXPECT().FindById(uint64(1)).Return(&game.Game{ID: 1, Title: "Test Game"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	context.SetPath("/games/:id")
	context.SetParamNames("id")
	context.SetParamValues("1")

	err := suite.handler.GetGame(context)
	suite.NoError(err)

	gameResponse := testing2.UnmarshalBody[game.Game](recorder)
	suite.Equal(http.StatusOK, recorder.Code)
	suite.Equal(game.Game{ID: 1, Title: "Test Game"}, gameResponse)
}

func (suite *HandlerTestSuite) TestFindByIdWithError() {
	suite.repo.EXPECT().FindById(uint64(1)).Return(nil, errors.New("some error"))

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	context.SetPath("/games/:id")
	context.SetParamNames("id")
	context.SetParamValues("1")

	err := suite.handler.GetGame(context)
	suite.NoError(err)

	errorDto := testing2.UnmarshalBody[dto.Error](recorder)
	suite.Equal(http.StatusInternalServerError, recorder.Code)
	suite.Equal("error loading game: some error", errorDto.Message)
}

// CreateGame

func (suite *HandlerTestSuite) TestCreateGame() {
	suite.repo.EXPECT().Create(mock.Anything).RunAndReturn(func(g *game.Game) error {
		g.ID = 1
		return nil
	}).Once()

	req, _ := http.NewRequest(
		http.MethodPost,
		"/games",
		testing2.MarshalBody(dto.Game{Title: "Test Game"}),
	)
	req.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	err := suite.handler.CreateGame(context)
	gameCreated := testing2.UnmarshalBody[game.Game](recorder)

	suite.NoError(err)
	suite.Equal(http.StatusCreated, recorder.Code)
	suite.Equal("Test Game", gameCreated.Title)
	suite.Equal(1, int(gameCreated.ID))
	suite.Equal(game.Game{ID: 1, Title: "Test Game"}, gameCreated)
}

func (suite *HandlerTestSuite) TestCreateGameWithBindingError() {
	req, _ := http.NewRequest(
		http.MethodPost,
		"/games",
		testing2.MarshalBody("{invalid_json}"),
	)
	req.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	err := suite.handler.CreateGame(context)
	errorDto := testing2.UnmarshalBody[dto.Error](recorder)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, recorder.Code)
	suite.Contains(errorDto.Message, "error binding game")
}

func (suite *HandlerTestSuite) TestCreateGameWithServiceError() {
	suite.repo.EXPECT().Create(mock.Anything).Return(errors.New("some error")).Once()

	req, _ := http.NewRequest(
		http.MethodPost,
		"/games",
		testing2.MarshalBody(dto.Game{Title: "Test Game"}),
	)
	req.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	context := suite.echo.NewContext(req, recorder)
	err := suite.handler.CreateGame(context)
	errorDto := testing2.UnmarshalBody[dto.Error](recorder)

	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, recorder.Code)
	suite.Contains(errorDto.Message, "error creating game")
}

func TestGamesHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
