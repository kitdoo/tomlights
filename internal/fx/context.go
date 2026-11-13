package fx

import (
	"context"
	"time"
)

// defaultInitTimeout caps initialization work so a hung dependency does not
// block startup forever.
const defaultInitTimeout = 30 * time.Second

// initContext creates a context for initialization operations. The caller
// MUST defer the returned cancel.
func initContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultInitTimeout)
}
