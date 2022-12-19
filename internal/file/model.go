package file

import (
	"io"
)

type File struct {
	Name        string `json:"name"`
	ShortID     string `json:"short_id"`
	ContentType string `json:"content_type"`
	io.Reader   `json:"-"`
}
