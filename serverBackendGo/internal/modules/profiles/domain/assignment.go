package domain

import "time"

// ProfileTreeAssignment links a published version to a device-tree folder.
type ProfileTreeAssignment struct {
	AssignmentID     int       `json:"assignmentId"`
	TreeNodeID       int       `json:"treeNodeId"`
	TreeNodeName     string    `json:"treeNodeName"`
	TreePath         string    `json:"treePath,omitempty"`
	ProfileVersionID int       `json:"profileVersionId"`
	VersionNumber    int       `json:"versionNumber"`
	DeviceCount      int       `json:"deviceCount"`
	CreatedAt        time.Time `json:"createdAt"`
}

// AssignmentImpact is preview before PUT assignment.
type AssignmentImpact struct {
	DeviceCount           int    `json:"deviceCount"`
	RequiresConfirmDialog bool   `json:"requiresConfirmDialog"`
	FolderName            string `json:"folderName"`
}

// PutAssignmentRequest is PUT /profiles/:id/assignments body.
type PutAssignmentRequest struct {
	TreeNodeID        int  `json:"treeNodeId"`
	ProfileVersionID  int  `json:"profileVersionId"`
	ConfirmImpact     bool `json:"confirmImpact"`
}

// PutAssignmentResult is PUT response data.
type PutAssignmentResult struct {
	ProfileTreeAssignment
	AffectedDevices int `json:"affectedDevices"`
}
