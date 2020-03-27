package persistent

import (
	"context"

	"github.com/iliyanmotovski/raytracer/backend"
)

// checkCtx check whether the context has been canceled or its deadline has expired
func checkCtx(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return backend.ErrContextCancelled
	case context.DeadlineExceeded:
		return backend.ErrContextExpired
	default:
		return nil
	}
}
