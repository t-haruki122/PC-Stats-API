package main

import (
	"embed"
	"io/fs"
)

//go:embed all:web
var webFS embed.FS

// GetWebFS returns the embedded web filesystem
func GetWebFS() fs.FS {
	// Strip "web/" prefix from paths
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		panic(err)
	}
	return sub
}
