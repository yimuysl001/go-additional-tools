package web

import (
	"embed"
)

//go:embed page/*
var staticFiles embed.FS
