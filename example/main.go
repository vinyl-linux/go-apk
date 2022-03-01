package main

import (
	"log"

	"github.com/vinyl-linux/go-apk"
)

var (
	baseURL = apk.BaseURL("edge", "community", "x86_64")
	client  = apk.Client{
		BaseURL: baseURL,
	}
)

func main() {
	// download index file
	apk.BaseDir = "tmp/"
	fn, err := client.DownloadIndex()
	if err != nil {
		panic(err)
	}

	// parse index
	a, err := apk.NewAPKIndex(fn)
	if err != nil {
		panic(err)
	}

	// print packages
	for idx, pkg := range a.Packages {
		log.Printf("%d\t\t%s (%s)",
			idx, pkg.Name, pkg.Description,
		)
	}

	// download an arbitrary package
	downloaded := a.Packages[321]
	p, err := client.DownloadPackage(downloaded)
	if err != nil {
		panic(err)
	}

	log.Printf("downloaded %s to %s", downloaded.Name, p)
}
