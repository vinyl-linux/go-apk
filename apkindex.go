package apk

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

// APKIndex represents the data stored in an APKINDEX.tar.gz file
type APKIndex struct {
	Description string
	Packages    Packages
}

// NewAPKIndex takes a filename containing an APKINDEX.tar.gz file,
// un-tars it, and processes the relevant files 'APKINDEX', and 'DESCRIPTION'
func NewAPKIndex(fn string) (a APKIndex, err error) {
	a.Packages, err = ReadPackages(fn)

	return
}

// Packages holds a series of packages, as described in an APKINDEX file
type Packages []*Package

// ReadPackages takes an APKINDEX files and parses the individual files contained
// therein
func ReadPackages(fn string) (p Packages, err error) {
	var (
		sb  strings.Builder
		pkg *Package
	)

	p = make(Packages, 0)

	f, err := os.Open(fn)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			pkg, err = DecodePackage(sb.String())
			if err != nil {
				return
			}

			p = append(p, pkg)
			sb.Reset()

			continue
		}

		sb.WriteString(line)
		sb.WriteString("\n")
	}

	// decode final
	pkg, err = DecodePackage(sb.String())
	if err != nil {
		return
	}

	p = append(p, pkg)

	err = scanner.Err()

	return
}

// Package represents a package from the Index file
// Keys and names and stuff are taken from
// https://github.com/alpinelinux/apk-tools/blob/da8d83338b7daa2bbd2e940f94d6515d69568800/src/apk_adb.c#L56-L76
type Package struct {
	ID            string           `apk:"C"`
	Name          string           `apk:"P"`
	Version       *version.Version `apk:"V"`
	Description   string           `apk:"T"`
	URL           string           `apk:"U"`
	InstalledSize int              `apk:"I"`
	FileSize      int              `apk:"S"`
	Licence       string           `apk:"L"`
	Arch          string           `apk:"A"`
	Dependencies  []string         `apk:"D"`
	InstallIf     string           `apk:"i"`
	Provides      []string         `apk:"p"`
	Origin        string           `apk:"o"`
	Maintainer    string           `apk:"m"`
	BuildTime     time.Time        `apk:"t"`
	RepoCommit    string           `apk:"c"`
	Replaces      string           `apk:"r"`
	Priority      int              `apk:"k"`
}

// DecodePackage takes a block of text and tries to decode a
// package from it
func DecodePackage(in string) (*Package, error) {
	p := &Package{}

	kvs := kv(in)

	ev := reflect.Indirect(reflect.ValueOf(p))
	et := ev.Type()

	for i := 0; i < et.NumField(); i++ {
		field, val := et.Field(i), ev.Field(i)

		key, ok := field.Tag.Lookup("apk")
		if !ok {
			continue
		}

		setVal, ok := kvs[key]
		if !ok {
			continue
		}

		switch tt := field.Type.String(); tt {
		case "string":
			val.SetString(setVal)

		case "*version.Version", "*Version":
			ver, err := version.NewVersion(setVal)
			if err != nil {
				return p, err
			}

			val.Set(reflect.ValueOf(ver))

		case "[]string":
			val.Set(reflect.ValueOf(strings.Fields(setVal)))

		case "time.Time", "Time":
			i, err := strconv.Atoi(setVal)
			if err != nil {
				return p, err
			}

			val.Set(reflect.ValueOf(time.Unix(int64(i), 0)))

		case "int":
			i, err := strconv.Atoi(setVal)
			if err != nil {
				return p, err
			}

			val.SetInt(int64(i))

		default:
			continue
		}
	}

	return p, nil
}

func kv(in string) (out map[string]string) {
	out = make(map[string]string)

	for _, line := range strings.Split(in, "\n") {
		elems := strings.SplitN(line, ":", 2)
		if elems == nil {
			continue
		}

		if len(elems) != 2 {
			continue
		}

		out[elems[0]] = elems[1]
	}

	return
}
