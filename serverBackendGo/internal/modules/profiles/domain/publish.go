package domain

// PublishImpactAssignment is one folder row in publish impact preview (020).
type PublishImpactAssignment struct {
	AssignmentID         int    `json:"assignmentId"`
	TreeNodeID           int    `json:"treeNodeId"`
	TreeNodeName         string `json:"treeNodeName"`
	CurrentVersionNumber int    `json:"currentVersionNumber"`
	DeviceCount          int    `json:"deviceCount"`
}
