// Package stamp applies text watermarks to PDF pages using pdfcpu.
package stamp

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"

	"github.com/anujkumar-df/pdfmark/internal/errs"
)

// NewTextWatermark builds a pdfcpu Watermark for the given text with the
// library's fixed default style: centered, diagonal, semi-transparent gray.
func NewTextWatermark(text string) *model.Watermark {
	wm := model.DefaultWatermarkConfig()
	wm.Mode = model.WMText
	wm.TextString = text
	wm.TextLines = []string{text}
	wm.OnTop = false
	wm.Pos = types.Center
	wm.Diagonal = model.DiagonalLLToUR
	wm.UserRotOrDiagonal = true
	wm.Opacity = 0.3
	wm.FontName = "Helvetica"
	wm.FontSize = 48
	wm.Scale = 1.0
	wm.ScaleAbs = false
	wm.Color = color.Gray
	wm.FillColor = color.Gray
	wm.StrokeColor = color.Gray
	return wm
}

// Apply reads the PDF from rs, stamps pages according to instructions
// (page number -> watermark text), and writes the result to w.
// The caller must have already validated that all page numbers are in range.
func Apply(rs io.ReadSeeker, w io.Writer, instructions map[int]string) error {
	if len(instructions) == 0 {
		_, err := io.Copy(w, rs)
		return err
	}

	wmMap := make(map[int]*model.Watermark, len(instructions))
	for page, text := range instructions {
		wmMap[page] = NewTextWatermark(text)
	}

	conf := model.NewDefaultConfiguration()
	return api.AddWatermarksMap(rs, w, wmMap, conf)
}

// PageCount returns the number of pages in the PDF behind rs.
// It resets the reader position to the start after counting.
func PageCount(rs io.ReadSeeker) (int, error) {
	conf := model.NewDefaultConfiguration()
	n, err := api.PageCount(rs, conf)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errs.ErrInvalidPDF, err)
	}
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		return 0, fmt.Errorf("seeking PDF: %w", err)
	}
	return n, nil
}

// ValidatePages checks that every page referenced in instructions exists
// within a PDF of totalPages pages.
func ValidatePages(instructions map[int]string, totalPages int) error {
	for page := range instructions {
		if page > totalPages {
			return fmt.Errorf("%w: page %d, PDF has %d pages", errs.ErrPageOutOfRange, page, totalPages)
		}
	}
	return nil
}

// BufferReader reads all of r into a bytes.Reader so pdfcpu can seek.
func BufferReader(r io.Reader) (*bytes.Reader, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading PDF input: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: empty input", errs.ErrInvalidPDF)
	}
	return bytes.NewReader(data), nil
}
