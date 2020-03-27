package backend

import "errors"

var (
	ErrContextExpired   = errors.New("context deadline exceeded")
	ErrContextCancelled = errors.New("context was canceled")
)
