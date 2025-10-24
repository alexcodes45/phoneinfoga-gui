package app

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yourorg/phoneinfoga-desktop/internal/scan"
	"github.com/yourorg/phoneinfoga-desktop/pkg/uiapi"
)

type App struct {
	Version string

	orchestrator *scan.Orchestrator
	api          *uiapi.API
}

func New(version string) *App {
	processor := func(ctx context.Context, job scan.Job) error {
		// Placeholder worker: in lieu of PhoneInfoga integration, simulate latency.
		select {
		case <-time.After(150 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	orchestrator := scan.New(processor)
	return &App{
		Version:      version,
		orchestrator: orchestrator,
		api:          uiapi.New(orchestrator),
	}
}

func (a *App) Startup(ctx context.Context) error {
	// TODO: wire config, DB, logger, DI
	log.Info().Str("version", a.Version).Msg("App starting")
	fmt.Println("App starting, version:", a.Version)
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.orchestrator.Shutdown(ctx)
}

// API returns the bridge used by Wails to expose backend features to the UI.
func (a *App) API() *uiapi.API {
	return a.api
}
