package apk

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	APKIndexFile    = "APKINDEX"
	DescriptionFile = "DESCRIPTION"
)

// APKIndex represents the data stored in an APKINDEX.tar.gz file
type APKIndex struct {
	Description []byte
	Packages    Packages
}

// NewAPKIndex takes a filename containing an APKINDEX.tar.gz file,
// un-tars it, and processes the relevant files 'APKINDEX', and 'DESCRIPTION'
func NewAPKIndex(fn string) (a APKIndex, err error) {
	// create tmpdir
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return
	}

	// extract fn to tmpdir
	err = extract(fn, dir)
	if err != nil {
		return
	}

	// read tmpdir/DESCRIPTION into a.Description
	a.Description, err = ioutil.ReadFile(filepath.Join(dir, DescriptionFile)) // #nosec: G304
	if err != nil {
		return
	}

	// read tmpdir/APKINDEX into packages
	a.Packages, err = ReadPackages(filepath.Join(dir, APKIndexFile)) // #nosec: G304

	return
}

func extract(fn, dir string) (err error) {
	f, err := os.Open(fn) // #nosec: G304
	if err != nil {
		return
	}

	decompressor, err := gzip.NewReader(f)
	if err != nil {
		return
	}

	defer decompressor.Close()

	tr := tar.NewReader(decompressor)
	return decompressLoop(tr, dir)
}

func decompressLoop(tr *tar.Reader, dest string) (err error) {
	var (
		header *tar.Header
		target string
	)

	for {
		header, err = tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target, err = sanitizeArchivePath(dest, header.Name)
		if err != nil {
			return
		}

		// ensure the parent directory exists in situations where a single file,
		// with missing toplevel dir, is compressed.
		//
		// This seemingly weird case exists in some tarballs in our test cases, so
		// it can definitely happen. I'd love to understand more /shrug
		err = os.MkdirAll(filepath.Dir(target), 0750)
		if err != nil {
			return
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0750); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode)) // #nosec: G304
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			err = f.Close()
			if err != nil {
				return err
			}
		}
	}
}

// sanitize archive file pathing from "G305: Zip Slip vulnerability"
//
// this function is taken almost verbatim from: https://github.com/securego/gosec/issues/324#issuecomment-935927967
func sanitizeArchivePath(d, t string) (v string, err error) {
	v = filepath.Join(d, t)
	if strings.HasPrefix(v, filepath.Clean(d)) {
		return v, nil
	}

	return "", fmt.Errorf("%s: %s", "content filepath is tainted", t)
}
