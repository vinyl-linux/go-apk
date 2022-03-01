package apk

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

var (
	APKIndexTarball = "APKINDEX.tar.gz"
	BaseDir         = "/var/cache/vinyl/vin/apk"
)

// Client connects to the Alpine package CDN in order to downloan
// indices, packages
type Client struct {
	BaseURL string
}

// DownloadIndex downloads an index file to a temporary directory,
// ready for parsing and processing, usually with apk.NewAPKIndex()
func (c Client) DownloadIndex() (path string, err error) {
	u := fmt.Sprintf("%s/%s", c.BaseURL, APKIndexTarball)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return
	}

	return download(u, dir)
}

// DownloadPackage accepts a package, as parsed from an APKIndex file, and
// stores it in the vinyl cache directory
func (c Client) DownloadPackage(pkg *Package) (path string, err error) {
	u := fmt.Sprintf("%s/%s-%s.apk", c.BaseURL, pkg.Name, pkg.Version)

	dir := filepath.Join(BaseDir, pkg.Name, pkg.Version)
	err = os.MkdirAll(dir, 0750)
	if err != nil {
		return
	}

	return download(u, dir)
}

func download(u string, dir string) (path string, err error) {
	u1, err := url.Parse(u)
	if err != nil {
		return
	}

	path = filepath.Join(dir, filepath.Base(u1.Path))
	out, err := os.Create(path)
	if err != nil {
		return
	}

	defer out.Close()

	resp, err := http.Get(u1.String())
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("downloading %s: expected status 200, received %s", u1, resp.Status)

		return
	}

	_, err = io.Copy(out, resp.Body)

	return
}
