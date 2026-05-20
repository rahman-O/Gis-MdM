package http

// LoginRequest mirrors legacy UserCredentials JSON.
// Password: raw text (Swagger/UI tools) or MD5 uppercase hex (production web/mobile clients).
type LoginRequest struct {
	Login    string `json:"login" example:"admin"`
	Password string `json:"password" example:"admin"`
}

// JWTResultDTO is the JWT login response body.
type JWTResultDTO struct {
	IDToken string `json:"id_token"`
}
