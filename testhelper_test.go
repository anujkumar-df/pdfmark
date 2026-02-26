package pdfmark

import (
	"io"
	"testing"

	"github.com/anujkumar-df/pdfmark/internal/testutil"
)

func createTestPDF(t *testing.T, pages int) []byte {
	t.Helper()
	return testutil.CreateTestPDF(t, pages)
}

func assertValidPDF(t *testing.T, data []byte) {
	t.Helper()
	testutil.AssertValidPDF(t, data)
}

func assertPageCount(t *testing.T, data []byte, expected int) {
	t.Helper()
	testutil.AssertPageCount(t, data, expected)
}

func csvString(lines ...string) io.Reader {
	return testutil.CSVString(lines...)
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }
