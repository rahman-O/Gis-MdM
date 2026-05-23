package domain

// PendingSignup is an in-progress customer registration.
type PendingSignup struct {
	ID         int64
	Email      string
	SignupTime int64
	Language   string
	Token      string
}

// SignupComplete is POST /public/signup/complete body.
type SignupComplete struct {
	Token       string
	Email       string
	Name        string
	FirstName   string
	LastName    string
	Company     string
	Description string
	PasswordMD5 string
}
