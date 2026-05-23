package port

import "context"

// Repository persists per-user hint history.
type Repository interface {
	GetHistory(ctx context.Context, userID int64) ([]string, error)
	MarkShown(ctx context.Context, userID int64, hintKey string) error
	Enable(ctx context.Context, userID int64) error
	Disable(ctx context.Context, userID int64) error
}
