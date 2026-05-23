package application

// AuthFailure indicates invalid credentials (maps to envelope ERROR).
type AuthFailure struct{}

func (AuthFailure) Error() string { return "authentication failed" }

// Unauthorized indicates JWT login failure (HTTP 401).
type Unauthorized struct{}

func (Unauthorized) Error() string { return "unauthorized" }

// BadRequest indicates missing credentials.
type BadRequest struct{}

func (BadRequest) Error() string { return "bad request" }
