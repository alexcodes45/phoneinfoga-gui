package app

import (
	"context"
	"fmt"
)

type App struct {
	Version string
}

func New(version string) *App {
	return &App{Version: version}
}

func (a *App) Startup(ctx context.Context) error {
	// TODO: wire config, DB, logger, DI
	fmt.Println("App starting, version:", a.Version)
	return nil
}
