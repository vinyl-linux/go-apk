README.md: doc.go
	goreadme -import-path github.com/vinyl-linux/go-apk -title go-apk -badge-godoc -badge-goreportcard -skip-sub-packages > $@
