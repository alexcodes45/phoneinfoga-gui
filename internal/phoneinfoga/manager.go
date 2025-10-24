package phoneinfoga

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Mode int

const (
	ModeServe Mode = iota
	ModeCLI
)

var (
	ErrBinaryMissing = errors.New("phoneinfoga binary not found")
	ErrServeTimeout  = errors.New("serve mode did not become ready")
	ErrInvalidNumber = errors.New("invalid phone number")
)

type Info struct {
	Version string
	Commit  string
}

type RequestOpts struct {
	ProxyURL string
	Timeout  time.Duration
	Headers  map[string]string
}

type RawResult []byte

type Client interface {
	Health(ctx context.Context) error
	Scan(ctx context.Context, number string, opts RequestOpts) (RawResult, error)
	Info(ctx context.Context) (Info, error)
}

type ServeOpts struct {
	Port       int
	BinaryPath string
	Env        []string
}

type Manager interface {
	EnsureServe(ctx context.Context, opts ServeOpts) (Client, error)
	ScanCLI(ctx context.Context, number string, opts RequestOpts) (RawResult, Info, error)
	DetectBinary(ctx context.Context) (string, error)
}

// ParseMode converts a textual mode representation to the enum value. Defaults to ModeServe.
func ParseMode(input string) Mode {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "cli":
		return ModeCLI
	case "serve":
		return ModeServe
	default:
		return ModeServe
	}
}

// String renders the mode to a human-readable string.
func (m Mode) String() string {
	switch m {
	case ModeCLI:
		return "cli"
	case ModeServe:
		fallthrough
	default:
		return "serve"
	}
}
