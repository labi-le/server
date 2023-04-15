package storage

import (
	"context"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"mime/multipart"
	"strings"
)

var (
	ErrInvalidFile          = errors.New("invalid file")
	ErrContentTypeAssertion = errors.New("invalid content type")
)

type Service interface {
	Add(ctx context.Context, rf RequestFile) (string, error)
	Get(ctx context.Context, hash string) (File, error)
}

type service struct {
	store Store
}

func NewService(c Store) Service {
	return &service{
		store: c,
	}
}

func (s *service) Add(_ context.Context, rf RequestFile) (string, error) {
	return rf.ShortID, s.store.Set(rf.ShortID, rf)
}

func (s *service) Get(_ context.Context, k string) (File, error) {
	var f File
	found, err := s.store.Get(k, &f)
	if err != nil {
		return f, err
	}

	if !found || f.Name == "" {
		return f, ErrFileNotFound
	}

	return f, nil
}

func getContentType(mp multipart.File) (*mimetype.MIME, error) {
	defer mp.Seek(0, io.SeekStart) //nolint:errcheck // dn

	return mimetype.DetectReader(mp)
}

const (
	alphabet    = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"
	alphabetLen = len(alphabet)
)

func Short(id int) string {
	var digits []int
	for id > 0 {
		digits = append(digits, id%alphabetLen)
		id /= alphabetLen
	}

	// reverse
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	var b strings.Builder
	for _, digit := range digits {
		b.WriteString(string(alphabet[digit]))
	}

	return b.String()
}
