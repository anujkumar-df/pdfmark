package pdfmark

import (
	"bytes"
	"errors"
	"os"
	"sync"
	"testing"
)

func TestWatermark_EndToEnd(t *testing.T) {
	pdf := createTestPDF(t, 5)

	csv := csvString(
		"page,watermark_text",
		"1,CONFIDENTIAL",
		"3,DRAFT",
		"5,INTERNAL USE ONLY",
	)

	var out bytes.Buffer
	err := Watermark(nopWriteCloser{&out}, bytes.NewReader(pdf), csv)
	if err != nil {
		t.Fatalf("Watermark: %v", err)
	}

	assertValidPDF(t, out.Bytes())
	assertPageCount(t, out.Bytes(), 5)

	if out.Len() <= len(pdf) {
		t.Error("output PDF should be larger than input after watermarking")
	}
}

func TestWatermark_NoWatermarks(t *testing.T) {
	pdf := createTestPDF(t, 3)
	csv := csvString("page,watermark_text")

	var out bytes.Buffer
	err := Watermark(nopWriteCloser{&out}, bytes.NewReader(pdf), csv)
	if err != nil {
		t.Fatalf("Watermark: %v", err)
	}

	assertValidPDF(t, out.Bytes())
	assertPageCount(t, out.Bytes(), 3)
}

func TestWatermark_AllPages(t *testing.T) {
	pdf := createTestPDF(t, 3)
	csv := csvString(
		"page,watermark_text",
		"1,PAGE ONE",
		"2,PAGE TWO",
		"3,PAGE THREE",
	)

	var out bytes.Buffer
	err := Watermark(nopWriteCloser{&out}, bytes.NewReader(pdf), csv)
	if err != nil {
		t.Fatalf("Watermark: %v", err)
	}

	assertValidPDF(t, out.Bytes())
	assertPageCount(t, out.Bytes(), 3)
}

func TestWatermark_PageOutOfRange(t *testing.T) {
	pdf := createTestPDF(t, 3)
	csv := csvString(
		"page,watermark_text",
		"1,OK",
		"10,BAD",
	)

	var out bytes.Buffer
	err := Watermark(nopWriteCloser{&out}, bytes.NewReader(pdf), csv)
	if !errors.Is(err, ErrPageOutOfRange) {
		t.Errorf("got error %v, want ErrPageOutOfRange", err)
	}
}

func TestWatermark_InvalidPDF(t *testing.T) {
	csv := csvString(
		"page,watermark_text",
		"1,TEST",
	)

	var out bytes.Buffer
	err := Watermark(nopWriteCloser{&out}, bytes.NewReader([]byte("not a pdf")), csv)
	if !errors.Is(err, ErrInvalidPDF) {
		t.Errorf("got error %v, want ErrInvalidPDF", err)
	}
}

func TestWatermark_EmptyPDF(t *testing.T) {
	csv := csvString(
		"page,watermark_text",
		"1,TEST",
	)

	var out bytes.Buffer
	err := Watermark(nopWriteCloser{&out}, bytes.NewReader(nil), csv)
	if !errors.Is(err, ErrInvalidPDF) {
		t.Errorf("got error %v, want ErrInvalidPDF", err)
	}
}

func TestWatermark_BadCSV(t *testing.T) {
	pdf := createTestPDF(t, 3)
	csv := csvString(
		"page,watermark_text",
		"abc,BAD",
	)

	var out bytes.Buffer
	err := Watermark(nopWriteCloser{&out}, bytes.NewReader(pdf), csv)
	if !errors.Is(err, ErrMalformedCSV) {
		t.Errorf("got error %v, want ErrMalformedCSV", err)
	}
}

func TestWatermark_DuplicatePageCSV(t *testing.T) {
	pdf := createTestPDF(t, 3)
	csv := csvString(
		"page,watermark_text",
		"1,FIRST",
		"1,SECOND",
	)

	var out bytes.Buffer
	err := Watermark(nopWriteCloser{&out}, bytes.NewReader(pdf), csv)
	if !errors.Is(err, ErrDuplicatePage) {
		t.Errorf("got error %v, want ErrDuplicatePage", err)
	}
}

func TestWatermark_WithFileFixtures(t *testing.T) {
	pdfData := createTestPDF(t, 5)
	csvData, err := os.ReadFile("testdata/valid.csv")
	if err != nil {
		t.Fatalf("reading CSV fixture: %v", err)
	}

	var out bytes.Buffer
	err = Watermark(nopWriteCloser{&out}, bytes.NewReader(pdfData), bytes.NewReader(csvData))
	if err != nil {
		t.Fatalf("Watermark: %v", err)
	}

	assertValidPDF(t, out.Bytes())
	assertPageCount(t, out.Bytes(), 5)
}

func TestWatermark_Concurrent(t *testing.T) {
	pdf := createTestPDF(t, 5)

	const goroutines = 10
	var wg sync.WaitGroup
	errs := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			csv := csvString(
				"page,watermark_text",
				"1,GOROUTINE",
				"3,CONCURRENT",
			)

			var out bytes.Buffer
			if err := Watermark(nopWriteCloser{&out}, bytes.NewReader(pdf), csv); err != nil {
				errs <- err
				return
			}

			assertValidPDF(t, out.Bytes())
			assertPageCount(t, out.Bytes(), 5)
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent error: %v", err)
	}
}
