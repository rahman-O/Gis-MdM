package auth

import "context"

// UserLookup loads users for JWT validation.
type UserLookup interface {
	LookupByLogin(ctx context.Context, login string) (*Principal, error)
}
