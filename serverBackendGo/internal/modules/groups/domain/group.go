package domain

// Group is a tenant-scoped device group.
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// LookupItem mirrors com.hmdm.rest.json.LookupItem.
type LookupItem struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}
