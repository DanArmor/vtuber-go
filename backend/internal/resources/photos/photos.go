package photos

import "embed"

const PhotosDir = "files"

//go:embed files/*
var PhotosFS embed.FS
