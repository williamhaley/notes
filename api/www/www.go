package www

import (
	"embed"
	"io/fs"
)

// content is our static web server content.
//
//go:embed public_html
var content embed.FS

func GetContent() fs.FS {
	content, err := fs.Sub(content, "public_html")
	if err != nil {
		panic(err)
	}

	return content
}
