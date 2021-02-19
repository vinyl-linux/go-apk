package apk

import (
	"fmt"
)

var (
	base = "http://dl-cdn.alpinelinux.org/alpine"
)

// BaseURL returns a URL which can be used to download packages, APKINDEX files
// from the alpine CDN
//
// This function provides a convenient way of deriving official looking URLs.
// 3rd party repos, or whatever, wont be supported here
func BaseURL(repository, packages, arch string) string {
	return fmt.Sprintf("%s/%s/%s/%s", base, repository, packages, arch)
}
