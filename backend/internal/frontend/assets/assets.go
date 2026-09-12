//go:build init

package assets

import _ "embed"

var assets = make(map[string][]byte)

func init() {
	assets[NotFoundImage] = notFoundImageBytes
	assets[FaviconImage] = faviconImageBytes
}

//go:embed not_found.png
var notFoundImageBytes []byte

//go:embed favicon.jpg
var faviconImageBytes []byte
