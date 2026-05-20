package domain

// AuthOptions is returned by GET /rest/public/auth/options.
type AuthOptions struct {
	Signup    bool    `json:"signup"`
	Recover   bool    `json:"recover"`
	PublicKey *string `json:"publicKey,omitempty"`
}
