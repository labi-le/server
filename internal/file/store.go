package file

import (
	"fmt"
	"io"
	"labile-me-serv/pkg/filesystem"
)

var ErrInvalidArgument = fmt.Errorf("invalid argument")
var ErrFileExists = fmt.Errorf("file already exists")
var ErrFileNotFound = fmt.Errorf("file not found")

// Store is an abstraction for different key-value store implementations.
// A store must be able to store, retrieve and delete key-value pairs,
// with the key being a string and the value being any Go interface{}.
type Store interface {
	// Set stores the given value for the given key.
	// The implementation automatically marshalls the value.
	// The marshalling format depends on the implementation. It can be JSON, gob etc.
	// The key must not be "" and the value must not be nil.
	Set(k string, v interface{}) error
	// Get retrieves the value for the given key.
	// The implementation automatically unmarshalls the value.
	// The unmarshalling source depends on the implementation. It can be JSON, gob etc.
	// The automatic unmarshalling requires a pointer to an object of the correct type
	// being passed as parameter.
	// In case of a struct the Get method will populate the fields of the object
	// that the passed pointer points to with the values of the retrieved object's values.
	// If no value is found it returns (false, nil).
	// The key must not be "" and the pointer must not be nil.
	Get(k string, v interface{}) (found bool, err error)
	// Delete deletes the stored value for the given key.
	// Deleting a non-existing key-value pair does NOT lead to an error.
	// The key must not be "".
	Delete(k string) error
	// Close must be called when the work with the key-value store is done.
	// Most (if not all) implementations are meant to be used long-lived,
	// so only call Close() at the very end.
	// Depending on the store implementation it might do one or more of the following:
	// Make sure all pending updates make their way to disk,
	// finish open transactions,
	// close the file handle to an embedded DB,
	// close the connection to the DB server,
	// release any open resources,
	// etc.
	// Some implementation might not need the store to be closed,
	// but as long as you work with the gokv.Store interface you never know which implementation
	// is passed to your method, so you should always call it.
	Close() error
}

type store struct {
	kvStore Store
	fs      filesystem.Storage
}

func NewStore(kvStore Store, fs filesystem.Storage) Store {
	return &store{kvStore: kvStore, fs: fs}
}

func (s *store) Set(k string, v interface{}) error {
	casted, ok := v.(RequestFile)
	if !ok {
		return ErrInvalidArgument
	}

	// check exist in kv store
	found, _ := s.kvStore.Get(k, &File{})
	if found {
		return ErrFileExists
	}

	file, err := s.fs.Create(casted.Name)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, copyErr := io.Copy(file, casted); copyErr != nil {
		return copyErr
	}

	return s.kvStore.Set(k, casted)
}

func (s *store) Get(k string, v interface{}) (bool, error) {
	casted, ok := v.(*File)
	if !ok {
		return false, ErrInvalidArgument
	}

	found, getErr := s.kvStore.Get(k, casted)
	if getErr != nil {
		return false, getErr
	}

	if !found {
		return false, ErrFileNotFound
	}

	if _, foundErr := s.fs.Stat(casted.Name); foundErr != nil {
		return false, ErrFileNotFound
	}

	ff, openErr := s.fs.Open(casted.Name)
	if openErr != nil {
		return false, openErr
	}

	casted.Reader = ff

	return true, nil
}

func (s *store) Delete(k string) error {
	if err := s.fs.Remove(k); err != nil {
		return err
	}

	return s.kvStore.Delete(k)
}

func (s *store) Close() error {
	return s.kvStore.Close()
}

//func hash(r io.Reader) string {
//	var buf bytes.Buffer
//	io.Copy(&buf, r)
//
//	sha := sha256.New()
//	sha.Write(buf.Bytes())
//
//	fileHash := fmt.Sprintf("%x", sha.Sum(nil))
//	return fileHash
//}
