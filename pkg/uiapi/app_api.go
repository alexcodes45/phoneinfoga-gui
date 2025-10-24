package uiapi

import (
	"context"
	"fmt"

	"github.com/yourorg/phoneinfoga-desktop/pkg/dto"
)

type API struct {}

func New() *API { return &API{} }

func (a *API) Scan(ctx context.Context, req dto.ScanRequest) (dto.ScanResponse, error) {
	// TODO: validate number, enqueue job via scan orchestrator
	fmt.Println("Scan requested:", req.Number, "mode:", req.Mode)
	return dto.ScanResponse{ScanID: 1, Status: "queued"}, nil
}
