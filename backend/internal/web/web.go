package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var files embed.FS

func Dist() (fs.FS, error) {
	return fs.Sub(files, "dist")
}
