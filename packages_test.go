package apk

import (
	"testing"
)

func TestReadPackages(t *testing.T) {
	for _, test := range []struct {
		name        string
		fn          string
		expectCount int
		expectError bool
	}{
		{"missing file", "testdata/nonsuch", 0, true},
		{"valid apkindex", "testdata/valid-APKINDEX", 2, false},
		{"dodgy version", "testdata/dodgy-version-APKINDEX", 0, true},
		{"dodgy time", "testdata/dodgy-time-APKINDEX", 0, true},
		{"dodgy priority", "testdata/dodgy-priority-APKINDEX", 0, true},
		{"dodgy multi", "testdata/dodgy-multi-APKINDEX", 0, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := ReadPackages(test.fn)
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %+v", err)
			}

			if test.expectCount != len(got) {
				t.Errorf("expected %d, received %d", test.expectCount, len(got))
			}
		})
	}
}
