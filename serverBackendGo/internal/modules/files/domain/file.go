package domain

// UploadedFile mirrors Java UploadedFile for API payloads.
type UploadedFile struct {
	ID               *int   `json:"id,omitempty"`
	CustomerID       int    `json:"customerId,omitempty"`
	FilePath         string `json:"filePath,omitempty"`
	Description      string `json:"description,omitempty"`
	UploadTime       int64  `json:"uploadTime,omitempty"`
	DevicePath       string `json:"devicePath,omitempty"`
	External         bool   `json:"external,omitempty"`
	ExternalURL      string `json:"externalUrl,omitempty"`
	ReplaceVariables bool   `json:"replaceVariables,omitempty"`
	URL              string `json:"url,omitempty"`
	TmpPath          string `json:"tmpPath,omitempty"`
	FileName         string `json:"fileName,omitempty"`
	Subdir           string `json:"subdir,omitempty"`
}

// FileView is the list row shape for the Files UI.
type FileView struct {
	ID                   *int     `json:"id,omitempty"`
	FilePath             string   `json:"filePath,omitempty"`
	Description          string   `json:"description,omitempty"`
	URL                  string   `json:"url,omitempty"`
	Size                 int64    `json:"size,omitempty"`
	UploadTime           int64    `json:"uploadTime,omitempty"`
	DevicePath           string   `json:"devicePath,omitempty"`
	External             bool     `json:"external,omitempty"`
	ReplaceVariables     bool     `json:"replaceVariables,omitempty"`
	UsedByIcons          []string `json:"usedByIcons,omitempty"`
	UsedByConfigurations []string `json:"usedByConfigurations,omitempty"`
}

// APKFileDetails subset for upload UX (matches Java APKFileDetails / React ApkFileDetails).
type APKFileDetails struct {
	Pkg         string `json:"pkg,omitempty"`
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	VersionCode int    `json:"versionCode,omitempty"`
	Arch        string `json:"arch,omitempty"`
}

// FileUploadResult returned from multipart upload.
type FileUploadResult struct {
	ServerPath  string          `json:"serverPath,omitempty"`
	FileDetails *APKFileDetails `json:"fileDetails,omitempty"`
	Application any             `json:"application,omitempty"`
	Complete    *bool           `json:"complete,omitempty"`
	Exists      *bool           `json:"exists,omitempty"`
	Name        string          `json:"name,omitempty"`
}

// LimitResponse storage quota for tenant.
type LimitResponse struct {
	SizeUsed  int `json:"sizeUsed"`
	SizeLimit int `json:"sizeLimit"`
}

// FileConfigurationLink configuration ↔ file association.
type FileConfigurationLink struct {
	ID                  *int   `json:"id,omitempty"`
	CustomerID          int    `json:"customerId,omitempty"`
	ConfigurationID     int    `json:"configurationId"`
	ConfigurationName   string `json:"configurationName,omitempty"`
	FileID              int    `json:"fileId"`
	FileName            string `json:"fileName,omitempty"`
	Upload              bool   `json:"upload,omitempty"`
	Remove              bool   `json:"remove,omitempty"`
	Notify              bool   `json:"notify,omitempty"`
}

// LinkConfigurationsToFileRequest bulk update body.
type LinkConfigurationsToFileRequest struct {
	Configurations []FileConfigurationLink `json:"configurations"`
}

// CustomerMeta tenant file storage settings.
type CustomerMeta struct {
	ID        int
	FilesDir  string
	SizeLimit int
	Master    bool
}
