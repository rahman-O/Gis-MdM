package domain

// NameResponse rebranding payload for login screens.
type NameResponse struct {
	AppName     string `json:"appName"`
	VendorName  string `json:"vendorName"`
	VendorLink  string `json:"vendorLink"`
	SignupLink  string `json:"signupLink"`
	TermsLink   string `json:"termsLink"`
}

// UploadAppRequest JSON in multipart field `app`.
type UploadAppRequest struct {
	DeviceID        string `json:"deviceId"`
	Hash            string `json:"hash"`
	Name            string `json:"name"`
	Pkg             string `json:"pkg"`
	Version         string `json:"version"`
	LocalPath       string `json:"localPath"`
	FileName        string `json:"fileName"`
	ShowIcon        bool   `json:"showIcon"`
	UseKiosk        bool   `json:"useKiosk"`
	RunAfterInstall bool   `json:"runAfterInstall"`
	RunAtBoot       bool   `json:"runAtBoot"`
	System          bool   `json:"system"`
}

// DeviceRef minimal device row for public upload.
type DeviceRef struct {
	ID         int
	CustomerID int
}
