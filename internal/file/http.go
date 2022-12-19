package file

import (
	"io"
)

type RequestFile struct {
	ShortID     string `json:"short_id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	io.Reader
}
