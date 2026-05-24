package application

import (
	"strings"
	"testing"
)

// wouldCreateCycle mirrors reparent validation in tree_repo.
func wouldCreateCycle(nodePath, newParentPath string) bool {
	return strings.HasPrefix(newParentPath, nodePath)
}

func TestWouldCreateCycle(t *testing.T) {
	tests := []struct {
		nodePath, newParentPath string
		want                    bool
	}{
		{"/1/2/", "/1/2/5/", true},
		{"/1/2/", "/1/3/", false},
		{"/1/", "/1/2/", true},
	}
	for _, tc := range tests {
		if got := wouldCreateCycle(tc.nodePath, tc.newParentPath); got != tc.want {
			t.Fatalf("cycle(%q, %q) = %v, want %v", tc.nodePath, tc.newParentPath, got, tc.want)
		}
	}
}
