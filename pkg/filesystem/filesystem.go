package filesystem

import (
	"github.com/spf13/afero"
	"os"
)

type Filesystem struct {
	afero.Fs
}

func New(path string) Storage {
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(err)
	}
	return &Filesystem{afero.NewBasePathFs(afero.NewOsFs(), path)}
}

func NewMemFS(path string) Storage {
	return &Filesystem{afero.NewBasePathFs(afero.NewMemMapFs(), path)}
}

func (f *Filesystem) Create(name string) (File, error) {
	return f.Fs.Create(name)
}

func (f *Filesystem) Open(name string) (File, error) {
	return f.Fs.Open(name)
}

func (f *Filesystem) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return f.Fs.OpenFile(name, flag, perm)
}
