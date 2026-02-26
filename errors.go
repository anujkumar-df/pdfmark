package pdfmark

import "github.com/anujkumar-df/pdfmark/internal/errs"

// Sentinel errors returned by Watermark.
var (
	ErrPageOutOfRange = errs.ErrPageOutOfRange
	ErrInvalidPage    = errs.ErrInvalidPage
	ErrDuplicatePage  = errs.ErrDuplicatePage
	ErrMalformedCSV   = errs.ErrMalformedCSV
	ErrInvalidPDF     = errs.ErrInvalidPDF
	ErrEmptyCSV       = errs.ErrEmptyCSV
)
