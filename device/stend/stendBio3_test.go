package stend

import (
	"testing"

)

func TestPlacesNumbers(t *testing.T) {
	for i,x := range Places {
		if i + 1 != len(x) {
			t.Errorf("%d: %d, %v", i + 1, len(x), x)
		}
	}
}