package pdfmark

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func createTestPDF(t *testing.T, pages int) []byte {
	t.Helper()
	if pages < 1 {
		t.Fatal("createTestPDF: pages must be >= 1")
	}

	conf := model.NewDefaultConfiguration()
	dim := &types.Dim{Width: 595.276, Height: 841.890}

	xRefTable, err := pdfcpu.CreateXRefTableWithRootDict()
	if err != nil {
		t.Fatalf("creating xref table: %v", err)
	}

	rootDict, err := xRefTable.Catalog()
	if err != nil {
		t.Fatalf("getting root dict: %v", err)
	}

	mediaBox := types.RectForDim(dim.Width, dim.Height)

	pagesDict := types.Dict(map[string]types.Object{
		"Type":     types.Name("Pages"),
		"Count":    types.Integer(pages),
		"MediaBox": mediaBox.Array(),
	})

	pagesIndRef, err := xRefTable.IndRefForNewObject(pagesDict)
	if err != nil {
		t.Fatalf("creating pages dict: %v", err)
	}

	kids := make(types.Array, 0, pages)
	for i := 0; i < pages; i++ {
		pageDict := types.Dict(map[string]types.Object{
			"Type":   types.Name("Page"),
			"Parent": *pagesIndRef,
		})
		pageIndRef, err := xRefTable.IndRefForNewObject(pageDict)
		if err != nil {
			t.Fatalf("creating page %d: %v", i+1, err)
		}
		kids = append(kids, *pageIndRef)
	}

	pagesDict.Insert("Kids", kids)
	rootDict.Insert("Pages", *pagesIndRef)

	ctx := pdfcpu.CreateContext(xRefTable, conf)

	var buf bytes.Buffer
	if err := api.WriteContext(ctx, &buf); err != nil {
		t.Fatalf("writing PDF: %v", err)
	}

	return buf.Bytes()
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func csvString(lines ...string) io.Reader {
	var buf bytes.Buffer
	for i, line := range lines {
		buf.WriteString(line)
		if i < len(lines)-1 {
			buf.WriteByte('\n')
		}
	}
	return &buf
}

func assertPageCount(t *testing.T, data []byte, expected int) {
	t.Helper()
	rs := bytes.NewReader(data)
	conf := model.NewDefaultConfiguration()
	n, err := api.PageCount(rs, conf)
	if err != nil {
		t.Fatalf("counting pages: %v", err)
	}
	if n != expected {
		t.Errorf("page count = %d, want %d", n, expected)
	}
}

func assertValidPDF(t *testing.T, data []byte) {
	t.Helper()
	rs := bytes.NewReader(data)
	conf := model.NewDefaultConfiguration()
	if err := api.Validate(rs, conf); err != nil {
		t.Fatalf("invalid PDF: %v", err)
	}
}

type errWriter struct {
	n int
}

func (e *errWriter) Write(p []byte) (int, error) {
	if e.n <= 0 {
		return 0, fmt.Errorf("simulated write error")
	}
	if len(p) > e.n {
		e.n = 0
		return len(p), fmt.Errorf("simulated write error")
	}
	e.n -= len(p)
	return len(p), nil
}

func (e *errWriter) Close() error { return nil }
