// +build !race

package uid

import (
	"testing"
)

func TestBasicUniqueness(t *testing.T) {
	n := 10000000
	m := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		n := Next()
		if _, ok := m[n]; ok {
			t.Fatalf("Duplicate NUID found: %v\n", n)
		}
		m[n] = struct{}{}
	}
}
