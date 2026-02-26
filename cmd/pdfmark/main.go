package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anujkumar-df/pdfmark"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func main() {
	pdfPath := flag.String("pdf", "", "path to input PDF")
	csvPath := flag.String("csv", "", "path to CSV watermark file")
	outPath := flag.String("out", "output.pdf", "path to output PDF")
	demo := flag.Bool("demo", false, "run a self-contained demo (ignores -pdf and -csv)")
	flag.Parse()

	if *demo || (*pdfPath == "" && *csvPath == "") {
		runDemo(*outPath)
		return
	}

	if *pdfPath == "" || *csvPath == "" {
		fmt.Fprintln(os.Stderr, "usage: pdfmark -pdf input.pdf -csv watermarks.csv [-out output.pdf]")
		fmt.Fprintln(os.Stderr, "       pdfmark -demo [-out output.pdf]")
		os.Exit(1)
	}

	pdfFile, err := os.Open(*pdfPath)
	if err != nil {
		log.Fatalf("opening PDF: %v", err)
	}
	defer pdfFile.Close()

	csvFile, err := os.Open(*csvPath)
	if err != nil {
		log.Fatalf("opening CSV: %v", err)
	}
	defer csvFile.Close()

	outFile, err := os.Create(*outPath)
	if err != nil {
		log.Fatalf("creating output: %v", err)
	}
	defer outFile.Close()

	if err := pdfmark.Watermark(outFile, pdfFile, csvFile); err != nil {
		log.Fatalf("watermarking failed: %v", err)
	}

	fmt.Printf("Done. Watermarked PDF written to %s\n", *outPath)
}

func runDemo(outPath string) {
	fmt.Println("Running demo mode...")
	fmt.Println()

	pdfData := generateDemoPDF(5)
	fmt.Printf("  Generated a %d-page demo PDF (%d bytes)\n", 5, len(pdfData))

	csvContent := strings.Join([]string{
		"page,watermark_text",
		"1,CONFIDENTIAL",
		"2,DRAFT",
		"4,INTERNAL USE ONLY",
	}, "\n")
	fmt.Printf("  CSV watermarks:\n")
	fmt.Printf("    Page 1 -> CONFIDENTIAL\n")
	fmt.Printf("    Page 2 -> DRAFT\n")
	fmt.Printf("    Page 4 -> INTERNAL USE ONLY\n")
	fmt.Printf("    Pages 3 and 5 -> (no watermark)\n")
	fmt.Println()

	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("creating output: %v", err)
	}
	defer outFile.Close()

	err = pdfmark.Watermark(outFile, bytes.NewReader(pdfData), strings.NewReader(csvContent))
	if err != nil {
		log.Fatalf("watermarking failed: %v", err)
	}

	fmt.Printf("Done. Open %s to see the watermarked PDF.\n", outPath)
}

func generateDemoPDF(pages int) []byte {
	conf := model.NewDefaultConfiguration()
	dim := &types.Dim{Width: 595.276, Height: 841.890}

	xRefTable, err := pdfcpu.CreateXRefTableWithRootDict()
	if err != nil {
		log.Fatalf("creating PDF: %v", err)
	}

	rootDict, err := xRefTable.Catalog()
	if err != nil {
		log.Fatalf("creating PDF: %v", err)
	}

	mediaBox := types.RectForDim(dim.Width, dim.Height)
	pagesDict := types.Dict(map[string]types.Object{
		"Type":     types.Name("Pages"),
		"Count":    types.Integer(pages),
		"MediaBox": mediaBox.Array(),
	})

	pagesIndRef, err := xRefTable.IndRefForNewObject(pagesDict)
	if err != nil {
		log.Fatalf("creating PDF: %v", err)
	}

	kids := make(types.Array, 0, pages)
	for i := 0; i < pages; i++ {
		pageDict := types.Dict(map[string]types.Object{
			"Type":   types.Name("Page"),
			"Parent": *pagesIndRef,
		})
		ref, err := xRefTable.IndRefForNewObject(pageDict)
		if err != nil {
			log.Fatalf("creating PDF: %v", err)
		}
		kids = append(kids, *ref)
	}

	pagesDict.Insert("Kids", kids)
	rootDict.Insert("Pages", *pagesIndRef)

	ctx := pdfcpu.CreateContext(xRefTable, conf)

	var buf bytes.Buffer
	if err := api.WriteContext(ctx, &buf); err != nil {
		log.Fatalf("writing PDF: %v", err)
	}

	return buf.Bytes()
}
