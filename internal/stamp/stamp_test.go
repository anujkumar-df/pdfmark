package stamp

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/anujkumar-df/pdfmark/internal/errs"
)

func TestApply_SinglePage(t *testing.T) {
	pdf := createTestPDF(t, 3)
	rs := bytes.NewReader(pdf)

	instructions := map[int]string{2: "CONFIDENTIAL"}

	var buf bytes.Buffer
	if err := Apply(rs, &buf, instructions); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	assertValidPDF(t, buf.Bytes())
	assertPageCount(t, buf.Bytes(), 3)
}

func TestApply_MultiplePages(t *testing.T) {
	pdf := createTestPDF(t, 5)
	rs := bytes.NewReader(pdf)

	instructions := map[int]string{
		1: "DRAFT",
		3: "CONFIDENTIAL",
		5: "INTERNAL",
	}

	var buf bytes.Buffer
	if err := Apply(rs, &buf, instructions); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	assertValidPDF(t, buf.Bytes())
	assertPageCount(t, buf.Bytes(), 5)
}

func TestApply_NoInstructions(t *testing.T) {
	pdf := createTestPDF(t, 2)
	rs := bytes.NewReader(pdf)

	var buf bytes.Buffer
	if err := Apply(rs, &buf, map[int]string{}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	assertValidPDF(t, buf.Bytes())
	assertPageCount(t, buf.Bytes(), 2)
}

func TestValidatePages_InRange(t *testing.T) {
	instructions := map[int]string{1: "A", 3: "B", 5: "C"}
	if err := ValidatePages(instructions, 5); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidatePages_OutOfRange(t *testing.T) {
	instructions := map[int]string{1: "A", 6: "B"}
	err := ValidatePages(instructions, 5)
	if !errors.Is(err, errs.ErrPageOutOfRange) {
		t.Errorf("got error %v, want ErrPageOutOfRange", err)
	}
}

func TestPageCount_Valid(t *testing.T) {
	for _, n := range []int{1, 3, 5} {
		pdf := createTestPDF(t, n)
		rs := bytes.NewReader(pdf)
		got, err := PageCount(rs)
		if err != nil {
			t.Fatalf("PageCount(%d): %v", n, err)
		}
		if got != n {
			t.Errorf("PageCount(%d) = %d, want %d", n, got, n)
		}
	}
}

func TestPageCount_InvalidPDF(t *testing.T) {
	rs := bytes.NewReader([]byte("not a pdf"))
	_, err := PageCount(rs)
	if !errors.Is(err, errs.ErrInvalidPDF) {
		t.Errorf("got error %v, want ErrInvalidPDF", err)
	}
}

func TestBufferReader_Empty(t *testing.T) {
	_, err := BufferReader(bytes.NewReader(nil))
	if !errors.Is(err, errs.ErrInvalidPDF) {
		t.Errorf("got error %v, want ErrInvalidPDF", err)
	}
}

func TestBufferReader_Valid(t *testing.T) {
	data := []byte("some pdf bytes")
	rs, err := BufferReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, _ := io.ReadAll(rs)
	if !bytes.Equal(out, data) {
		t.Errorf("got %q, want %q", out, data)
	}
}
