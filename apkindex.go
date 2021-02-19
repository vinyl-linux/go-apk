package apk

// APKIndex represents the data stored in an APKINDEX.tar.gz file
type APKIndex struct {
	Description string
	Packages    Packages
}

// NewAPKIndex takes a filename containing an APKINDEX.tar.gz file,
// un-tars it, and processes the relevant files 'APKINDEX', and 'DESCRIPTION'
func NewAPKIndex(fn string) (a APKIndex, err error) {
	// download

	a.Packages, err = ReadPackages(fn)

	return
}
