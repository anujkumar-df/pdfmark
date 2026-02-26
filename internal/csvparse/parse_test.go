package csvparse

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/anujkumar-df/pdfmark/internal/errs"
	"github.com/anujkumar-df/pdfmark/internal/testutil"
)

func csvString(lines ...string) io.Reader {
	return testutil.CSVString(lines...)
}

func TestParse_Valid(t *testing.T) {
	r := csvString(
		"page,watermark_text",
		"1,CONFIDENTIAL",
		"3,DRAFT",
		"5,INTERNAL USE ONLY",
	)

	m, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 3 {
		t.Fatalf("got %d entries, want 3", len(m))
	}
	if m[1] != "CONFIDENTIAL" {
		t.Errorf("page 1 = %q, want %q", m[1], "CONFIDENTIAL")
	}
	if m[3] != "DRAFT" {
		t.Errorf("page 3 = %q, want %q", m[3], "DRAFT")
	}
	if m[5] != "INTERNAL USE ONLY" {
		t.Errorf("page 5 = %q, want %q", m[5], "INTERNAL USE ONLY")
	}
}

func TestParse_EmptyBody(t *testing.T) {
	r := csvString("page,watermark_text")
	m, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("got %d entries, want 0", len(m))
	}
}

func TestParse_EmptyInput(t *testing.T) {
	r := strings.NewReader("")
	_, err := Parse(r)
	if !errors.Is(err, errs.ErrEmptyCSV) {
		t.Errorf("got error %v, want ErrEmptyCSV", err)
	}
}

func TestParse_DuplicatePage(t *testing.T) {
	r := csvString(
		"page,watermark_text",
		"1,CONFIDENTIAL",
		"1,DRAFT",
	)
	_, err := Parse(r)
	if !errors.Is(err, errs.ErrDuplicatePage) {
		t.Errorf("got error %v, want ErrDuplicatePage", err)
	}
}

func TestParse_InvalidPageNumber(t *testing.T) {
	r := csvString(
		"page,watermark_text",
		"abc,CONFIDENTIAL",
	)
	_, err := Parse(r)
	if !errors.Is(err, errs.ErrMalformedCSV) {
		t.Errorf("got error %v, want ErrMalformedCSV", err)
	}
}

func TestParse_NegativePage(t *testing.T) {
	r := csvString(
		"page,watermark_text",
		"-1,CONFIDENTIAL",
	)
	_, err := Parse(r)
	if !errors.Is(err, errs.ErrInvalidPage) {
		t.Errorf("got error %v, want ErrInvalidPage", err)
	}
}

func TestParse_ZeroPage(t *testing.T) {
	r := csvString(
		"page,watermark_text",
		"0,CONFIDENTIAL",
	)
	_, err := Parse(r)
	if !errors.Is(err, errs.ErrInvalidPage) {
		t.Errorf("got error %v, want ErrInvalidPage", err)
	}
}

func TestParse_MissingColumns(t *testing.T) {
	r := csvString(
		"page,watermark_text",
		"1",
	)
	_, err := Parse(r)
	if !errors.Is(err, errs.ErrMalformedCSV) {
		t.Errorf("got error %v, want ErrMalformedCSV", err)
	}
}

func TestParse_EmptyWatermarkText(t *testing.T) {
	r := csvString(
		"page,watermark_text",
		"1,",
	)
	_, err := Parse(r)
	if !errors.Is(err, errs.ErrMalformedCSV) {
		t.Errorf("got error %v, want ErrMalformedCSV", err)
	}
}

func TestParse_HeaderOnlyOneColumn(t *testing.T) {
	r := csvString("page")
	_, err := Parse(r)
	if !errors.Is(err, errs.ErrMalformedCSV) {
		t.Errorf("got error %v, want ErrMalformedCSV", err)
	}
}

func TestParse_ExtraColumns(t *testing.T) {
	r := csvString(
		"page,watermark_text,extra",
		"1,CONFIDENTIAL,ignored",
	)
	m, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m[1] != "CONFIDENTIAL" {
		t.Errorf("page 1 = %q, want %q", m[1], "CONFIDENTIAL")
	}
}

func TestParse_WhitespaceTrimming(t *testing.T) {
	r := csvString(
		"page, watermark_text",
		" 2 , DRAFT ",
	)
	m, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m[2] != "DRAFT" {
		t.Errorf("page 2 = %q, want %q", m[2], "DRAFT")
	}
}
