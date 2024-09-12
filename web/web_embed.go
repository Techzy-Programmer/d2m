package web

import (
	"embed"
)

//go:embed panel/dist/*
var EmbeddedFiles embed.FS
