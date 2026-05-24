package domain

import "errors"

// Sentinel errors for HTTP mapping.
var (
	ErrNotFound        = errors.New("error.notfound.object")
	ErrDuplicateName   = errors.New("error.device_tree.duplicate_name")
	ErrInvalidParent   = errors.New("error.device_tree.invalid_parent")
	ErrCycle           = errors.New("error.device_tree.cycle")
	ErrTargetRequired  = errors.New("error.device_tree.target_required")
	ErrCannotDeleteRoot = errors.New("error.device_tree.cannot_delete_root")
)

const RootFolderName = "All devices"

// TreeNode is a tenant-scoped folder in the device tree.
type TreeNode struct {
	ID          int    `json:"id"`
	ParentID    *int   `json:"parentId"`
	Name        string `json:"name"`
	SortOrder   int    `json:"sortOrder"`
	Path        string `json:"path"`
	Depth       int    `json:"depth"`
	DeviceCount int    `json:"deviceCount"`
}

// TreeList is returned by GET /device-tree.
type TreeList struct {
	Nodes  []TreeNode `json:"nodes"`
	RootID int        `json:"rootId"`
}

// CreateNodeRequest is POST /device-tree/nodes body.
type CreateNodeRequest struct {
	ParentID  int    `json:"parentId"`
	Name      string `json:"name"`
	SortOrder int    `json:"sortOrder"`
}

// UpdateNodeRequest is PUT /device-tree/nodes/:id body.
type UpdateNodeRequest struct {
	Name      *string `json:"name"`
	SortOrder *int    `json:"sortOrder"`
	ParentID  *int    `json:"parentId"`
}

// DeleteNodeRequest is POST /device-tree/nodes/:id/delete body.
type DeleteNodeRequest struct {
	TargetNodeID int `json:"targetNodeId"`
}
