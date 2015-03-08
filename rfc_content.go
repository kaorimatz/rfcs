package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type DefaultRFCContentRepository struct {
	Fetcher    *RFCContentFetcher
	CacheStore *RFCContentCacheStore
}

func (r *DefaultRFCContentRepository) FindByNumber(number int) ([]byte, error) {
	if r.CacheStore != nil {
		if content, _ := r.CacheStore.Get(number); content != nil {
			return content, nil
		}
	}

	content, err := r.Fetcher.Fetch(number)
	if err != nil {
		return nil, err
	}

	if r.CacheStore != nil {
		r.CacheStore.Put(number, content)
	}

	return content, nil
}

func NewDefaultRFCContentRepository() *DefaultRFCContentRepository {
	repository := DefaultRFCContentRepository{
		Fetcher:    &RFCContentFetcher{FileFormat: RFCContentFileFormatASCII},
		CacheStore: &RFCContentCacheStore{},
	}

	return &repository
}

type RFCContentFetcher struct {
	FileFormat RFCContentFileFormat
}

func (f *RFCContentFetcher) Fetch(number int) ([]byte, error) {
	rfcURL, err := f.FileFormat.URLFor(number)
	if err != nil {
		return nil, err
	}

	response, err := http.Get(rfcURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

type RFCContentFileFormat int

const (
	RFCContentFileFormatASCII RFCContentFileFormat = iota
	RFCContentFileFormatPs
	RFCContentFileFormatPdf
)

func (f RFCContentFileFormat) URLFor(number int) (string, error) {
	switch f {
	case RFCContentFileFormatASCII:
		return fmt.Sprintf("http://www.rfc-editor.org/rfc/rfc%d.txt", number), nil
	case RFCContentFileFormatPs:
		return fmt.Sprintf("http://www.rfc-editor.org/rfc/rfc%d.ps", number), nil
	case RFCContentFileFormatPdf:
		return fmt.Sprintf("http://www.rfc-editor.org/rfc/rfc%d.pdf", number), nil
	}
	return "", fmt.Errorf("no URL available for file format: %v", f)
}

type RFCContentCacheStore struct {
	CacheDirectory string
}

func (s *RFCContentCacheStore) Put(number int, content []byte) error {
	cacheDir, err := s.cacheDirectory()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	cacheFile := filepath.Join(cacheDir, strconv.Itoa(number))

	return ioutil.WriteFile(cacheFile, content, 0644)
}

func (s *RFCContentCacheStore) Get(number int) ([]byte, error) {
	cacheDir, err := s.cacheDirectory()
	if err != nil {
		return nil, err
	}

	cacheFile := filepath.Join(cacheDir, strconv.Itoa(number))

	if _, err = os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, nil
	}

	return ioutil.ReadFile(cacheFile)
}

func (s *RFCContentCacheStore) cacheDirectory() (string, error) {
	if s.CacheDirectory != "" {
		return s.CacheDirectory, nil
	}

	if dir := GetUserCacheDirectory(); dir != "" {
		return dir, nil
	}

	return "", fmt.Errorf("cannot determine the cache directory")
}
