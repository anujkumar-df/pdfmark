// Package errs defines sentinel errors used throughout the pdfmark library.
package errs

import "errors"

var (
	ErrPageOutOfRange = errors.New("pdfmark: page number out of range")
	ErrInvalidPage    = errors.New("pdfmark: page number must be >= 1")
	ErrDuplicatePage  = errors.New("pdfmark: duplicate page number in CSV")
	ErrMalformedCSV   = errors.New("pdfmark: malformed CSV row")
	ErrInvalidPDF     = errors.New("pdfmark: invalid or corrupt PDF input")
	ErrEmptyCSV       = errors.New("pdfmark: CSV contains no header row")
)
