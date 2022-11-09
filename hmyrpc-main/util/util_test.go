package util

import (
	"reflect"
	"testing"
)

func TestIsZeroHash(t *testing.T) {
	zeroHash := "0x0000000000000000000000000000000000000000000000000000000000000000"
	nonZeroHash := "0x0000000000000000000000000000000000000000000000000000000000000001"
	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Fatalf(" %s mismatch: got %v want %v", f, got, want)
		}
	}
	// check
	check("zero hash", IsZeroHash(zeroHash), true)
	check("non zero hash", IsZeroHash(nonZeroHash), false)
}
