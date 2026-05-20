package application

// Logout is a no-op at application layer; session cleared in HTTP adapter.
func (s *Service) Logout() {}
