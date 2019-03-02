package models

// This is in a separate package to avoid circular import cycles.

// Part struct
type Part struct {
	MimeType string `json:"mimetype,omitempty"`
	Data     string `json:"data,omitempty"`
}
