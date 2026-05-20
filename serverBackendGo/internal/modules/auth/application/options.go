package application

import (
	"context"

	"github.com/gis-mdm/server-backend-go/internal/modules/auth/domain"
)

// Options returns login page options.
func (s *Service) Options(_ context.Context) domain.AuthOptions {
	opts := domain.AuthOptions{
		Signup:  s.email.IsConfigured() && s.customerSignup,
		Recover: s.email.IsConfigured(),
	}
	if s.transmitPassword && s.rsa != nil {
		if pk, err := s.rsa.PublicKeyBase64(); err == nil && pk != "" {
			opts.PublicKey = &pk
		}
	}
	return opts
}
