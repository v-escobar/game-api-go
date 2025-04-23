package interfaces

import (
	"github.com/v-escobar/game-api-go/internal/interfaces/common"
	"github.com/v-escobar/game-api-go/internal/interfaces/docs"
	"github.com/v-escobar/game-api-go/internal/interfaces/game"
)

func NewHandlers(
	gameHandler *game.Handler,
	documentationHandler *docs.DocumentationHandler,
) []common.ApplicationHandler {
	return []common.ApplicationHandler{
		gameHandler, documentationHandler,
	}
}
