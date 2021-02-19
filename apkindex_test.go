package apk

import (
	"testing"
)

func TestNewAPKIndex(t *testing.T) {
	for _, test := range []struct {
		name        string
		fn          string
		expectError bool
	}{
		{"happy path", "testdata/valid-APKINDEX.tar.gz", false},
		{"missing description", "testdata/missing-description-APKINDEX.tar.gz", true},
		{"missing file", "testdata/no-such-APKINDEX.tar.gz", true},
		{"non-trivial tarball", "testdata/non-trivial.tar.gz", false},
	} {
		t.Run(test.name, func(f *testing.T) {
			_, err := NewAPKIndex(test.fn)
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %+v", err)
			}
		})
	}
}
