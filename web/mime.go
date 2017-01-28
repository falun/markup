package web

import (
	"mime"
	"path"
)

type MimeTypeProvider interface {
	ForFile(path string) string
}

type provider struct {
	assetDir string
}

func NewMimeTypeProvider(assetDir string) MimeTypeProvider {
	return provider{assetDir}
}

func (p provider) ForFile(fp string) string {
	return mime.TypeByExtension(path.Ext(fp))
}
