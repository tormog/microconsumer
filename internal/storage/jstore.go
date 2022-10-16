package storage

import (
	"os"

	js "github.com/schollz/jsonstore"
)

func NewJSONStorage(fileStore string) Storage {
	if _, err := os.Stat(fileStore); err == nil {
		data, err := js.Open(fileStore)
		if err == nil {
			return &jsonStorage{
				fileStore: fileStore,
				jStore:    data,
			}
		}
	}
	return &jsonStorage{
		fileStore: fileStore,
		jStore:    new(js.JSONStore),
	}
}

func (s *jsonStorage) Close() error {
	// Saving will automatically gzip if .gz is provided
	if err := js.Save(s.jStore, s.fileStore); err != nil {
		return err
	}
	return nil
}

func (s *jsonStorage) PersistData(id, source string, value interface{}) error {
	err := s.jStore.Set(id, value)
	if err != nil {
		return err
	}
	return nil
}
