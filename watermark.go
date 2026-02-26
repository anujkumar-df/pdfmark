package pdfmark

import (
	"io"

	"github.com/anujkumar-df/pdfmark/internal/csvparse"
	"github.com/anujkumar-df/pdfmark/internal/stamp"
)

// Watermark reads a PDF from src, applies text watermarks according to the CSV
// data from csvData, and writes the resulting PDF to dst.
//
// The CSV must have a header row and at least two columns: page (1-indexed) and
// watermark_text. Pages not listed in the CSV are passed through unchanged.
//
// The caller is responsible for closing dst; this function only writes to it.
//
// Watermark is safe for concurrent use from multiple goroutines.
func Watermark(dst io.WriteCloser, src io.Reader, csvData io.Reader) error {
	instructions, err := csvparse.Parse(csvData)
	if err != nil {
		return err
	}

	rs, err := stamp.BufferReader(src)
	if err != nil {
		return err
	}

	totalPages, err := stamp.PageCount(rs)
	if err != nil {
		return err
	}

	if err := stamp.ValidatePages(instructions, totalPages); err != nil {
		return err
	}

	return stamp.Apply(rs, dst, instructions)
}
