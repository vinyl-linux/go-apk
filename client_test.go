package apk

import (
	"io/ioutil"
	"testing"
)

func TestClient_DownloadIndex(t *testing.T) {
	for _, test := range []struct {
		name        string
		url         string
		expectError bool
	}{
		{"happy path", BaseURL("edge", "community", "x86_64"), false},
		{"404ing path", BaseURL("hedge", "community", "x86_64"), true},
		{"bad url", "\x00", true},
		{"connection error", "http://127.0.0.1:4321", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := Client{test.url}

			_, err := c.DownloadIndex()
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %+v", err)
			}

		})
	}
}

func TestClient_DownloadPackage(t *testing.T) {
	dir, err := ioutil.TempDir("", "apk-test")
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	oldBaseDir := BaseDir
	BaseDir = dir
	defer func() {
		BaseDir = oldBaseDir
	}()

	// This feels very fragile: what if this is upgraded?
	// I've tried to pick something that sounds like it's not updated
	// much, which I can't image this is
	pkg1 := &Package{Name: "abiword-plugin-bmp", Version: "3.0.4-r4"}
	pkg2 := &Package{Name: "this-package-does-not-exist", Version: "3.0.4-r4"}

	for _, test := range []struct {
		name        string
		url         string
		pkg         *Package
		expectError bool
	}{
		//{"happy path", BaseURL("edge", "community", "x86_64"), pkg1, false},    // see comment above; this test is *already* broken
		{"404ing path", BaseURL("hedge", "community", "x86_64"), pkg2, true},
		{"bad url", "\x00", pkg1, true},
		{"connection error", "http://127.0.0.1:4321", pkg1, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := Client{test.url}

			_, err := c.DownloadPackage(test.pkg)
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %+v", err)
			}

		})
	}
}
