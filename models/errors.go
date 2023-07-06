package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrEmailTaken = errors.New("models: Email address is already in use")
	ErrNotFound   = errors.New("models: resource could not be found")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %w", fe.Issue)
}

func checkContentType(r io.ReadSeeker, allowedContents []string) error {
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	_, err = r.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("checking content type: %w", err)
	}

	contentType := http.DetectContentType(testBytes)
	for _, t := range allowedContents {
		if t == contentType {
			return nil
		}
	}

	return FileError{
		Issue: fmt.Sprintf("invalid content type: %v", contentType),
	}
}
