package models

// Part struct
type Part struct {
	Hex    string `json:"hex,omitempty"`
	UTF8   string `json:"utf8,omitempty"`
	BASE64 string `json:"base64,omitempty"`
	URI    string `json:"uri,omitempty"`
}
