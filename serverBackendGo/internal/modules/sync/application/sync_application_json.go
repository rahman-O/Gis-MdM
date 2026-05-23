package application

// BoolTrueOnly matches Java SyncApplication JSON: false booleans are omitted (NON_NULL).
func BoolTrueOnly(v bool) *bool {
	if v {
		b := true
		return &b
	}
	return nil
}
