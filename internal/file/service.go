package file

import (
	"context"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"mime/multipart"
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
