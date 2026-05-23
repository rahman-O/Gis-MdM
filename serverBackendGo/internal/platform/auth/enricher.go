package auth

import "context"

// PrincipalEnricher loads role and permission data onto a principal.
type PrincipalEnricher interface {
	EnrichPrincipal(ctx context.Context, p *Principal) error
}
