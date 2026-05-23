package email

import (
	"log/slog"
)

// Service sends transactional email when configured (dev: logs only).
type Service struct {
	configured bool
	log        *slog.Logger
}

func NewService(configured bool, log *slog.Logger) *Service {
	return &Service{configured: configured, log: log}
}

func (s *Service) IsConfigured() bool { return s.configured }

// Send logs the message when email is not configured (matches safe recover/signup behavior in dev).
func (s *Service) Send(to, subject, body string) bool {
	if s.configured {
		s.log.Info("email send (stub — configure SMTP in future)", "to", to, "subject", subject)
		return true
	}
	s.log.Debug("email skipped (EMAIL_CONFIGURED=false)", "to", to, "subject", subject)
	return true
}

func (s *Service) RecoveryBody(token string) string {
	return "Password reset token: " + token + "\nUse the link in your MDM UI: /password-reset/" + token
}

func (s *Service) VerifySignupBody(token string) string {
	return "Complete signup with token: " + token + "\n/open signup-complete/" + token
}
