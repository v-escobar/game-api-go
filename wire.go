//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/v-escobar/game-api-go/internal/application"
	"github.com/v-escobar/game-api-go/internal/infrastructure/db"
	"github.com/v-escobar/game-api-go/internal/infrastructure/persistence"
	"github.com/v-escobar/game-api-go/internal/interfaces"
	"github.com/v-escobar/game-api-go/internal/interfaces/docs"
	"github.com/v-escobar/game-api-go/internal/interfaces/game"
)

func initializeApp(dsn string) *App {
	wire.Build(
		NewApp,
		db.NewDB,
		game.NewGameHandler,
		application.NewGameService,
		persistence.NewGameRepository,
		interfaces.NewHandlers,
		docs.NewDocumentationHandler,
	)

	return &App{}
}
