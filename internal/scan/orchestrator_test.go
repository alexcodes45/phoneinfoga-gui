package scan

import (
	"context"
	"testing"
	"time"
)

func TestOrchestratorProcessesJobs(t *testing.T) {
	processor := func(ctx context.Context, job Job) error {
		select {
		case <-time.After(20 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	orch := New(processor, WithWorkers(1))
	defer orch.Shutdown(context.Background())

	ctx := context.Background()
	jobID, err := orch.Enqueue(ctx, Job{NumberE164: "+12025550187"})
	if err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	finished := waitForEvent(t, orch.Events(), jobID, EventFinished, 500*time.Millisecond)
	if !finished {
		t.Fatalf("did not observe finished event for %s", jobID)
	}

	stats := orch.Stats()
	if stats.Completed != 1 || stats.Running != 0 || stats.Pending != 0 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestOrchestratorCancel(t *testing.T) {
	processor := func(ctx context.Context, job Job) error {
		<-ctx.Done()
		return ctx.Err()
	}
	orch := New(processor, WithWorkers(1))
	defer orch.Shutdown(context.Background())

	ctx := context.Background()
	jobID, err := orch.Enqueue(ctx, Job{NumberE164: "+12025550187"})
	if err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	if ok := waitForEvent(t, orch.Events(), jobID, EventStarted, 500*time.Millisecond); !ok {
		t.Fatalf("did not observe started event for %s", jobID)
	}

	if err := orch.Cancel(jobID); err != nil {
		t.Fatalf("cancel: %v", err)
	}

	if ok := waitForEvent(t, orch.Events(), jobID, EventCanceled, 500*time.Millisecond); !ok {
		t.Fatalf("did not observe canceled event for %s", jobID)
	}

	stats := orch.Stats()
	if stats.Canceled != 1 {
		t.Fatalf("expected 1 canceled job, got %+v", stats)
	}
}

func waitForEvent(t *testing.T, events <-chan Event, jobID string, target EventType, timeout time.Duration) bool {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case evt, ok := <-events:
			if !ok {
				return false
			}
			if evt.JobID != jobID {
				continue
			}
			if evt.Type == target {
				return true
			}
		case <-timer.C:
			return false
		}
	}
}
