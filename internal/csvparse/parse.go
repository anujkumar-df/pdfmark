// Package csvparse reads watermark instructions from CSV data.
package csvparse

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/anujkumar-df/pdfmark/internal/errs"
)

// Parse reads a CSV from r with the expected format:
//
//	page,watermark_text
//	1,CONFIDENTIAL
//	3,DRAFT
//
// It returns a map from 1-indexed page number to watermark text.
// Errors are returned for missing headers, duplicate pages, non-integer page
// numbers, pages <= 0, and rows with fewer than 2 fields.
func Parse(r io.Reader) (map[int]string, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true

	header, err := cr.Read()
	if err == io.EOF {
		return nil, errs.ErrEmptyCSV
	}
	if err != nil {
		return nil, fmt.Errorf("%w: reading header: %v", errs.ErrMalformedCSV, err)
	}
	if len(header) < 2 {
		return nil, fmt.Errorf("%w: header must have at least 2 columns, got %d", errs.ErrMalformedCSV, len(header))
	}

	instructions := make(map[int]string)

	for line := 2; ; line++ {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%w: line %d: %v", errs.ErrMalformedCSV, line, err)
		}

		if len(record) < 2 {
			return nil, fmt.Errorf("%w: line %d: expected at least 2 fields, got %d", errs.ErrMalformedCSV, line, len(record))
		}

		pageStr := strings.TrimSpace(record[0])
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("%w: line %d: invalid page number %q", errs.ErrMalformedCSV, line, pageStr)
		}

		if page <= 0 {
			return nil, fmt.Errorf("%w: page %d on line %d", errs.ErrInvalidPage, page, line)
		}

		text := strings.TrimSpace(record[1])
		if text == "" {
			return nil, fmt.Errorf("%w: line %d: watermark text is empty", errs.ErrMalformedCSV, line)
		}

		if _, exists := instructions[page]; exists {
			return nil, fmt.Errorf("%w: page %d on line %d", errs.ErrDuplicatePage, page, line)
		}

		instructions[page] = text
	}

	return instructions, nil
}
