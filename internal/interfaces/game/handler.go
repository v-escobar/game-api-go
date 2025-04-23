package game

import (
	"errors"
	"fmt"
	"github.com/v-escobar/game-api-go/internal/application"
	applicationErrors "github.com/v-escobar/game-api-go/internal/application/errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/v-escobar/game-api-go/internal/interfaces/dto"
)

type Handler struct {
	service application.GameService
}

func NewGameHandler(service application.GameService) *Handler {
	return &Handler{service}
}

func (h *Handler) BindRoutes(e *echo.Echo) {
	g := e.Group("/games")
	{
		g.POST("", h.CreateGame)
		g.GET("", h.ListGames)
		g.GET("/:id", h.GetGame)
	}
}

// @Summary		Creates a new game
// @Description	Creates a new game
// @Accept			json
// @Produce		json
// @Router			/games [post]
func (h *Handler) CreateGame(c echo.Context) error {
	var game dto.Game
	if err := c.Bind(&game); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{Message: fmt.Sprint("error binding game: ", err.Error())})
	}

	if err := h.service.Create(&game); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error{Message: fmt.Sprint("error creating game: ", err.Error())})
	}

	return c.JSON(http.StatusCreated, game)
}

func (h *Handler) ListGames(c echo.Context) error {
	games, err := h.service.FindAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error{Message: fmt.Sprint("could not fetch games, Error: ", err.Error())})
	}

	return c.JSON(http.StatusOK, games)
}

func (h *Handler) GetGame(c echo.Context) error {
	gameIdString := c.Param("id")
	gameId, err := strconv.ParseUint(gameIdString, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{Message: fmt.Sprint("error parsing identifier: ", err.Error())})
	}

	persistedGame, err := h.service.FindById(gameId)

	if err != nil {
		switch {
		case errors.Is(err, applicationErrors.GameNotFoundError):
			return c.JSON(http.StatusNotFound, dto.Error{Message: fmt.Sprint("Game not found with id ", gameId)})
		default:
			return c.JSON(http.StatusInternalServerError, dto.Error{Message: fmt.Sprint("error loading game: ", err.Error())})
		}
	}

	return c.JSON(http.StatusOK, persistedGame)
}
