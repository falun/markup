package web

import (
	"mime"
	"path/filepath"
)

type MimeTypeProvider interface {
	ForFile(path string) string
}

type provider struct {
	assetDir     string
	extensionMap map[string]string
}

func NewMimeTypeProvider(assetDir string, extMap map[string]string) MimeTypeProvider {
	if extMap == nil {
		extMap = map[string]string{}
	}
	return provider{assetDir, extMap}
}

func (p provider) ForFile(fp string) string {
	ext := filepath.Ext(fp)
	if mt, ok := p.extensionMap[ext]; ok {
		return mt
	}

	return mime.TypeByExtension(ext)
}
