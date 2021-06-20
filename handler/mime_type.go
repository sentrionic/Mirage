package handler

var validImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// isAllowedImageType determines if image is among types defined
// in map of allowed images
func isAllowedImageType(mimeType string) bool {
	_, exists := validImageTypes[mimeType]

	return exists
}
