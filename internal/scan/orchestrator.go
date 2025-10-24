package scan

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog/log"
	"github.com/yourorg/phoneinfoga-desktop/internal/phoneinfoga"
)

// EventType represents the type of job lifecycle event emitted by the orchestrator.
type EventType string

const (
	EventStarted  EventType = "started"
	EventProgress EventType = "progress"
	EventFinished EventType = "finished"
	EventError    EventType = "error"
	EventCanceled EventType = "canceled"
)

var (
	// ErrQueueClosed is returned when attempting to enqueue after shutdown.
	ErrQueueClosed = errors.New("scan orchestrator has been shut down")
	jobSeq         uint64
)

// Job describes a PhoneInfoga scan request scheduled for execution.
type Job struct {
	ID            string
	NumberE164    string
	PreferredMode phoneinfoga.Mode
	Options       phoneinfoga.RequestOpts
	CaseID        *int
}

// Event is emitted on state changes for a job.
type Event struct {
	JobID   string
	Type    EventType
	Detail  string
	Percent float64
}

// Stats captures queue and execution counters.
type Stats struct {
	Pending   int
	Running   int
	Completed int
	Failed    int
	Canceled  int
}

// Processor executes the actual scan work for a job.
type Processor func(ctx context.Context, job Job) error

type queuedJob struct {
	ctx    context.Context
	job    Job
	cancel context.CancelFunc
}

type orchestratorConfig struct {
	workers   int
	queueSize int
}

// Option configures the orchestrator behaviour.
type Option func(*orchestratorConfig)

// WithWorkers sets the number of concurrent workers (default: 2).
func WithWorkers(n int) Option {
	return func(cfg *orchestratorConfig) {
		if n > 0 {
			cfg.workers = n
		}
	}
}

// WithQueueSize sets the size of the pending job buffer (default: 32).
func WithQueueSize(n int) Option {
	return func(cfg *orchestratorConfig) {
		if n > 0 {
			cfg.queueSize = n
		}
	}
}

// Orchestrator coordinates scan jobs, worker execution, and lifecycle events.
type Orchestrator struct {
	processor Processor

	queue   chan queuedJob
	events  chan Event
	cancels map[string]context.CancelFunc

	wg      sync.WaitGroup
	mu      sync.RWMutex
	statsMu sync.RWMutex
	closed  bool
	stats   Stats
}

// New builds an orchestrator using the supplied processor.
func New(processor Processor, opts ...Option) *Orchestrator {
	cfg := orchestratorConfig{workers: 2, queueSize: 32}
	for _, opt := range opts {
		opt(&cfg)
	}
	o := &Orchestrator{
		processor: processor,
		queue:     make(chan queuedJob, cfg.queueSize),
		events:    make(chan Event, cfg.queueSize),
		cancels:   make(map[string]context.CancelFunc),
	}
	o.wg.Add(cfg.workers)
	for i := 0; i < cfg.workers; i++ {
		go o.worker()
	}
	return o
}

// Events exposes the orchestrator event stream.
func (o *Orchestrator) Events() <-chan Event { return o.events }

// Enqueue schedules a job for execution.
func (o *Orchestrator) Enqueue(ctx context.Context, job Job) (string, error) {
	if job.NumberE164 == "" {
		return "", fmt.Errorf("NumberE164 is required")
	}
	if job.ID == "" {
		job.ID = o.nextID()
	}

	o.mu.RLock()
	closed := o.closed
	o.mu.RUnlock()
	if closed {
		return "", ErrQueueClosed
	}

	jobCtx, cancel := context.WithCancel(ctx)
	queued := queuedJob{ctx: jobCtx, job: job, cancel: cancel}

	select {
	case o.queue <- queued:
		o.storeCancel(job.ID, cancel)
		o.bumpStats(func(s *Stats) { s.Pending++ })
		return job.ID, nil
	case <-ctx.Done():
		cancel()
		return "", ctx.Err()
	}
}

// Cancel requests cancellation of a queued or running job.
func (o *Orchestrator) Cancel(jobID string) error {
	o.mu.RLock()
	closed := o.closed
	o.mu.RUnlock()
	if closed {
		return ErrQueueClosed
	}
	cancel := o.takeCancel(jobID)
	if cancel == nil {
		log.Warn().Str("job", jobID).Msg("cancel requested for unknown job")
		return fmt.Errorf("job %s not found", jobID)
	}
	cancel()
	return nil
}

// Stats returns a snapshot of queue metrics.
func (o *Orchestrator) Stats() Stats {
	o.statsMu.RLock()
	defer o.statsMu.RUnlock()
	return o.stats
}

// Shutdown drains the queue and waits for workers to exit.
func (o *Orchestrator) Shutdown(ctx context.Context) error {
	o.mu.Lock()
	if o.closed {
		o.mu.Unlock()
		return nil
	}
	o.closed = true
	close(o.queue)
	o.mu.Unlock()

	done := make(chan struct{})
	go func() {
		o.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		close(o.events)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (o *Orchestrator) worker() {
	defer o.wg.Done()
	for queued := range o.queue {
		job := queued.job
		o.bumpStats(func(s *Stats) {
			if s.Pending > 0 {
				s.Pending--
			}
			s.Running++
		})
		o.events <- Event{JobID: job.ID, Type: EventStarted}

		err := o.processor(queued.ctx, job)
		o.takeCancel(job.ID)

		select {
		case <-queued.ctx.Done():
			o.events <- Event{JobID: job.ID, Type: EventCanceled, Detail: queued.ctx.Err().Error()}
			o.bumpStats(func(s *Stats) {
				if s.Running > 0 {
					s.Running--
				}
				s.Canceled++
			})
		default:
			if err != nil {
				o.events <- Event{JobID: job.ID, Type: EventError, Detail: err.Error()}
				o.bumpStats(func(s *Stats) {
					if s.Running > 0 {
						s.Running--
					}
					s.Failed++
				})
			} else {
				o.events <- Event{JobID: job.ID, Type: EventFinished}
				o.bumpStats(func(s *Stats) {
					if s.Running > 0 {
						s.Running--
					}
					s.Completed++
				})
			}
		}
	}
}

func (o *Orchestrator) storeCancel(id string, cancel context.CancelFunc) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.cancels[id] = cancel
}

func (o *Orchestrator) takeCancel(id string) context.CancelFunc {
	o.mu.Lock()
	defer o.mu.Unlock()
	cancel, ok := o.cancels[id]
	if ok {
		delete(o.cancels, id)
	}
	return cancel
}

func (o *Orchestrator) bumpStats(mutator func(*Stats)) {
	o.statsMu.Lock()
	mutator(&o.stats)
	o.statsMu.Unlock()
}

func (o *Orchestrator) nextID() string {
	next := atomic.AddUint64(&jobSeq, 1)
	return fmt.Sprintf("job-%d", next)
}

// EmitProgress allows a processor to push progress updates to listeners.
func (o *Orchestrator) EmitProgress(jobID string, percent float64, detail string) {
	o.events <- Event{JobID: jobID, Type: EventProgress, Percent: percent, Detail: detail}
}
