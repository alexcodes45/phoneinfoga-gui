package scan

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yourorg/phoneinfoga-desktop/internal/phoneinfoga"
)

type Job struct {
	ID            string
	NumberE164    string
	PreferredMode phoneinfoga.Mode
	Options       phoneinfoga.RequestOpts
	CaseID        *int
}

type Event struct {
	JobID  string
	Type   string // started|progress|finished|error|canceled
	Detail string
}

type Orchestrator struct {
	evch chan Event
	wg   sync.WaitGroup
}

func New() *Orchestrator {
	return &Orchestrator{evch: make(chan Event, 64)}
}

func (o *Orchestrator) Events() <-chan Event { return o.evch }

func (o *Orchestrator) Enqueue(ctx context.Context, job Job) (string, error) {
	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		o.evch <- Event{JobID: job.ID, Type: "started"}
		// TODO: validate number, pick serve/cli, persist DB, etc.
		time.Sleep(500 * time.Millisecond)
		o.evch <- Event{JobID: job.ID, Type: "finished", Detail: "stub complete"}
	}()
	return job.ID, nil
}

func (o *Orchestrator) Cancel(ctx context.Context, jobID string) error {
	log.Info().Str("job", jobID).Msg("cancel requested (not yet implemented)")
	return nil
}
