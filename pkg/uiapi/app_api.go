package uiapi

import (
	"context"
	"fmt"
	"time"

	"github.com/yourorg/phoneinfoga-desktop/internal/phoneinfoga"
	"github.com/yourorg/phoneinfoga-desktop/internal/scan"
	"github.com/yourorg/phoneinfoga-desktop/pkg/dto"
)

// API exposes backend functionality to the Wails front-end.
type API struct {
	orchestrator  *scan.Orchestrator
	defaultRegion string
}

// New constructs the API bridge with the provided orchestrator.
func New(orchestrator *scan.Orchestrator) *API {
	region := "US"
	return &API{orchestrator: orchestrator, defaultRegion: region}
}

// Scan validates the number and enqueues a new job for processing.
func (a *API) Scan(ctx context.Context, req dto.ScanRequest) (dto.ScanResponse, error) {
	if a.orchestrator == nil {
		return dto.ScanResponse{}, fmt.Errorf("scan orchestrator not configured")
	}

	normalized, err := scan.NormalizeE164(req.Number, a.defaultRegion)
	if err != nil {
		return dto.ScanResponse{}, err
	}

	mode := phoneinfoga.ParseMode(req.Mode)
	opts := phoneinfoga.RequestOpts{
		ProxyURL: req.Proxy,
	}
	if req.TimeoutMs > 0 {
		opts.Timeout = time.Duration(req.TimeoutMs) * time.Millisecond
	}

	job := scan.Job{
		NumberE164:    normalized,
		PreferredMode: mode,
		Options:       opts,
		CaseID:        req.CaseID,
	}

	jobID, err := a.orchestrator.Enqueue(ctx, job)
	if err != nil {
		return dto.ScanResponse{}, err
	}

	return dto.ScanResponse{JobID: jobID, Status: "queued", Number: normalized}, nil
}
