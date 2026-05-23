package application

import (
	"github.com/gis-mdm/server-backend-go/internal/modules/auth/port"
	platformcrypto "github.com/gis-mdm/server-backend-go/internal/platform/crypto"
	platformemail "github.com/gis-mdm/server-backend-go/internal/platform/email"
	platformjwt "github.com/gis-mdm/server-backend-go/internal/platform/jwt"
)

// Service groups auth use cases.
type Service struct {
	repo             port.UserRepository
	jwt              *platformjwt.Provider
	email            *platformemail.Service
	rsa              *platformcrypto.RSAKeys
	transmitPassword bool
	customerSignup   bool
}

// NewService constructs the auth application service.
func NewService(
	repo port.UserRepository,
	jwt *platformjwt.Provider,
	email *platformemail.Service,
	rsa *platformcrypto.RSAKeys,
	transmitPassword, customerSignup bool,
) *Service {
	return &Service{
		repo:             repo,
		jwt:              jwt,
		email:            email,
		rsa:              rsa,
		transmitPassword: transmitPassword,
		customerSignup:   customerSignup,
	}
}

// Repo exposes repository for sibling modules (passwordreset, signup, twofactor, users).
func (s *Service) Repo() port.UserRepository { return s.repo }
