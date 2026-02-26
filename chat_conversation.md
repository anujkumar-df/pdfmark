# Golang PDF generation library requirements
_Exported on 2/26/2026 at 07:45:21 GMT+5:30 from Cursor (2.5.25)_

---

**User**

Sharing a requirement, do not write code yet, lets first understand what is the actual requirement

Deliver a working and Production Ready Golang lib and Test cases which can produce a PDF document with inputs as io.Reader for the PDF template, CSV as a Watermark data for desired pages as needed. The output should be written to io.WriteCloser.

---

**Cursor**

Great, let's break this down carefully before writing any code. I want to make sure I understand the full scope of what you need.

Here's my understanding so far:

**Core idea:** A Go library that takes an existing PDF (template), applies text watermarks to specified pages based on CSV data, and writes the resulting PDF to an output stream.

I have several clarifying questions that will significantly shape the design:

Good, that narrows things down significantly. A few more important questions:

Excellent. Now I have a complete picture. Let me put together the plan.

---

**User**

suggest a name for the library.

---

**Cursor**

Here are a few naming options — considering Go conventions (short, lowercase, no hyphens in package names):

---

**User**

https://github.com/anujkumar-df/pdfmark.git 

Configure this as origin

---

**Cursor**

Got it. Let me initialize git, set the remote, and update the plan with the module path `github.com/anujkumar-df/pdfmark`.

PowerShell syntax issue. Let me run them sequentially.

Now let me update the plan to reflect the confirmed library name and module path.

Done. Here's a summary of what's set up:

- **Git repo** initialized at `c:\Users\Admin\work\test-golang`
- **Remote origin** set to `https://github.com/anujkumar-df/pdfmark.git`
- **Plan updated** — module path is now `github.com/anujkumar-df/pdfmark`, package name is `pdfmark`

The library name **pdfmark** is clean, short, and idiomatic Go. When you're ready, confirm the plan and I'll start building it out.

---

**User**

Lets commit the plan to the repo first.

---

**Cursor**

I'll copy the plan into the repo as a design document (stripping the Cursor-specific frontmatter) and commit it.

---

**User**

Implement the plan as specified, it is attached for your reference. Do NOT edit the plan file itself.

To-do's from the plan have already been created. Do not create them again. Mark them as in_progress as you work, starting with the first one. Don't stop until you have completed all the to-dos.

---

**Cursor**

Starting with module initialization. Let me also check what Go version is available.

Good, pdfcpu v0.11.1 is installed. Let me quickly check its API surface before writing code.

`AddWatermarksMap` takes `map[int]*model.Watermark` directly — perfect match for our page-to-watermark mapping. Let me check how to construct a `model.Watermark`.

Now I have a clear picture of the pdfcpu API. Let me also check the `ParseTextWatermarkDetails` function signature to understand the `desc` parameter format.

Good. Now I have all the API knowledge I need. Let me build the library. I'll create the errors, CSV parser, stamp logic, and public API files.

Now the stamp logic.

Now the public API.

Let me verify the code compiles before moving on to test fixtures.

Clean compile. Now let me create the test data fixtures. I'll write a small Go helper to generate a multi-page test PDF using pdfcpu, and create the CSV files.

I need to generate a test PDF. Let me write a small Go program to create a 5-page PDF using pdfcpu, then create the CSV fixtures.

The API signature is different. Let me take a simpler approach — I'll write a test helper that generates the PDF programmatically, and write the CSV fixtures as plain files.

Let me check pdfcpu's import API — it can create PDFs from images or content descriptors.

Let me write a test helper that creates PDFs, and then the CSV fixtures.

Good. So the strategy is: create a 1-page context using `CreateContextWithXRefTable`, write it, then use `InsertPages` to add more blank pages. But actually, `CreateContextWithXRefTable` creates a context WITHOUT pages. Let me look at how to add the first page.

Let me try a cleaner approach using the exported functions.

Good, compiles clean. Let me verify the test helper actually works with a quick smoke test.

The demo page function needs content. Let me build the page tree directly without content.

I'll construct the page tree manually — a blank page doesn't need content streams.

Test PDF helper works. Now let me create the CSV fixture files.

Now let me delete the smoke test file and write the real unit tests and integration tests.

Now the stamp unit tests.

All 21 unit tests pass. Now for the integration and concurrency tests.

All integration and concurrency tests pass with the race detector enabled. Now let me run the full suite and verify everything.

All **31 tests pass** with the race detector. Let me also run `go vet` and check lints.

Clean across the board. Let me verify the final file structure.

All tasks complete. Here's a summary of what was built:

---

**Library: `github.com/anujkumar-df/pdfmark`**

**Source files (4):**
- `watermark.go` — Public API: single `Watermark(dst io.WriteCloser, src io.Reader, csvData io.Reader) error` function
- `csv.go` — CSV parser with strict validation (duplicate pages, invalid page numbers, malformed rows)
- `stamp.go` — pdfcpu integration: buffering, page counting, validation, and watermark application
- `errors.go` — Sentinel errors: `ErrPageOutOfRange`, `ErrInvalidPage`, `ErrDuplicatePage`, `ErrMalformedCSV`, `ErrInvalidPDF`, `ErrEmptyCSV`

**Test files (3) + test helpers (1):**
- `csv_test.go` — 12 unit tests covering all CSV parsing edge cases
- `stamp_test.go` — 9 unit tests for stamping, page counting, validation, and buffering
- `watermark_test.go` — 10 integration tests including end-to-end, error paths, file fixtures, and a concurrent goroutine safety test
- `testhelper_test.go` — Shared test utilities (PDF generation, assertions, writers)

**Test fixtures (5 CSV files):**
- `testdata/valid.csv`, `invalid_page.csv`, `duplicate_page.csv`, `empty.csv`, `malformed.csv`

**Results:**
- **31/31 tests passing**
- **Race detector clean** (`-race` flag)
- **`go vet` clean**
- **Zero linter errors**
- Concurrency-safe (stateless design, verified with 10 concurrent goroutines)

---

**User**

Arrange th src code in hierarchical folder Structure

---

**Cursor**

Let me first review the current structure and then reorganize it properly.

Here's the target structure I'll implement:

```
pdfmark/
├── go.mod
├── go.sum
├── plan.md
├── doc.go                            # package-level documentation
├── errors.go                         # re-exports sentinel errors for public API
├── watermark.go                      # public Watermark() function
├── watermark_test.go                 # integration + concurrency tests
├── testhelper_test.go                # shared test utilities (root)
├── internal/
│   ├── errs/
│   │   └── errs.go                   # sentinel error definitions (no deps)
│   ├── csvparse/
│   │   ├── parse.go                  # CSV parsing logic
│   │   └── parse_test.go            # CSV unit tests
│   └── stamp/
│       ├── stamp.go                  # pdfcpu watermark application
│       ├── stamp_test.go            # stamp unit tests
│       └── testhelper_test.go       # PDF generation helper for stamp tests
└── testdata/
    └── *.csv                         # test fixtures (unchanged)
```

The key design choice: errors live in `internal/errs` (no circular deps) and are re-exported by the root `errors.go` so the public API surface stays identical. Let me build it.

Now the stamp internal package.

Now the stamp test helper and stamp tests.

Now update the root package files.

Now rewrite `errors.go` to re-export from the internal package, and `watermark.go` to import the internal packages.

Now update the root-level test files. The integration tests and test helpers need to use the internal packages.

Now update the integration test — it needs to use the `stamp` internal package for `pageCount` in the concurrent test (or just use the helper). Let me update it to use the public API and test helpers only, which is cleaner for integration tests.

Now delete the old flat files that have been moved to internal packages.

Now let me verify it all compiles and passes.

All **31 tests** pass across **3 packages** with the race detector. Let me also run `go vet`.

Let me verify the final file tree.

Done. Here's the final hierarchical structure:

```
pdfmark/
├── go.mod / go.sum              # module: github.com/anujkumar-df/pdfmark
├── plan.md                      # design document
│
├── doc.go                       # package-level documentation
├── errors.go                    # public sentinel errors (re-exports from internal)
├── watermark.go                 # public API: Watermark(dst, src, csv)
├── watermark_test.go            # 10 integration + concurrency tests
├── testhelper_test.go           # shared test utilities for root package
│
├── internal/
│   ├── errs/
│   │   └── errs.go              # sentinel error definitions (zero dependencies)
│   ├── csvparse/
│   │   ├── parse.go             # CSV parsing logic (imports errs)
│   │   └── parse_test.go        # 12 CSV unit tests
│   └── stamp/
│       ├── stamp.go             # pdfcpu watermark application (imports errs)
│       ├── stamp_test.go        # 9 stamp unit tests
│       └── testhelper_test.go   # PDF generation helper for stamp tests
│
└── testdata/
    └── *.csv                    # 5 CSV test fixtures
```

**Key design decisions:**

- `internal/errs` holds all sentinel errors with zero project imports, breaking any potential circular dependency
- `internal/csvparse` owns CSV parsing, fully tested in isolation
- `internal/stamp` owns all pdfcpu interaction, fully tested in isolation
- Root `errors.go` re-exports the internal errors so the **public API surface is identical** -- callers still write `errors.Is(err, pdfmark.ErrInvalidPDF)`
- Root `watermark.go` is a thin orchestrator that wires the internal packages together
- Each package has its own test helpers, keeping tests self-contained

**31/31 tests pass** across 3 packages, race-detector clean, `go vet` clean.

---

**User**

Now, how can I use this library in a different golang based executable?

---

**Cursor**

Here's how you'd use `pdfmark` from another Go project:

```go
package main

import (
	"log"
	"os"

	"github.com/anujkumar-df/pdfmark"
)

func main() {
	// Open the PDF template
	pdfFile, err := os.Open("template.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer pdfFile.Close()

	// Open the CSV with watermark instructions
	csvFile, err := os.Open("watermarks.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	// Create the output file
	outFile, err := os.Create("output.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Apply watermarks — one call, done
	if err := pdfmark.Watermark(outFile, pdfFile, csvFile); err != nil {
		log.Fatal(err)
	}

	log.Println("watermarked PDF written to output.pdf")
}
```

**Setup in the consuming project:**

```bash
go get github.com/anujkumar-df/pdfmark
```

**Error handling example** (checking specific failure types):

```go
import "errors"

err := pdfmark.Watermark(outFile, pdfFile, csvFile)
if err != nil {
	switch {
	case errors.Is(err, pdfmark.ErrPageOutOfRange):
		log.Fatal("CSV references a page that doesn't exist in the PDF")
	case errors.Is(err, pdfmark.ErrDuplicatePage):
		log.Fatal("CSV has duplicate page entries")
	case errors.Is(err, pdfmark.ErrMalformedCSV):
		log.Fatal("CSV format is invalid")
	case errors.Is(err, pdfmark.ErrInvalidPDF):
		log.Fatal("PDF input is corrupt or invalid")
	default:
		log.Fatalf("unexpected error: %v", err)
	}
}
```

The key points:

- All three arguments are standard interfaces (`io.Reader`, `io.Reader`, `io.WriteCloser`), so they work with files, buffers, HTTP request/response bodies, S3 streams, etc.
- `os.File` satisfies `io.WriteCloser` natively, so you pass it directly
- The caller owns the lifecycle of all three handles (open and close them yourself)
- Since `Watermark` is goroutine-safe, you can call it concurrently from an HTTP handler serving multiple requests

---

**User**

Create a main.go in this repo to test it out.

---

**Cursor**

Works. Let me verify the output is a valid watermarked PDF.

Output PDF is 2,967 bytes (up from 687 bytes input), confirming the watermarks were applied. The `output.pdf` is ready to open.

Here's how to use it:

**Demo mode** (generates everything for you):
```bash
go run ./cmd/pdfmark -demo
```

**With your own files:**
```bash
go run ./cmd/pdfmark -pdf template.pdf -csv watermarks.csv -out result.pdf
```

**Flags:**
- `-pdf` -- path to your input PDF
- `-csv` -- path to your watermark CSV
- `-out` -- output path (defaults to `output.pdf`)
- `-demo` -- runs a self-contained demo, no files needed

The generated `output.pdf` is in the workspace root -- you can open it in any PDF viewer to see the diagonal watermarks on pages 1, 2, and 4.

---

**User**

"C:\Users\Admin\Downloads\NCERT Solutions for Class 6 Science.pdf"

Use this as input

---

**Cursor**

The PDF has only 4 pages, but `testdata/valid.csv` references page 5. The strict validation caught it correctly. Let me create a CSV that fits this 4-page PDF and re-run.

Done. The watermarked `output.pdf` has been written (148 KB). The NCERT PDF has 4 pages, and watermarks were applied as follows:

| Page | Watermark |
|------|-----------|
| 1 | CONFIDENTIAL |
| 2 | DRAFT |
| 3 | *(none)* |
| 4 | INTERNAL USE ONLY |

Open `output.pdf` in your PDF viewer to see the diagonal watermarks. Note the first run failing with `ErrPageOutOfRange` was the strict validation working exactly as designed -- it refused to proceed because `valid.csv` referenced page 5 on a 4-page document.

---

**User**

Does pdf have layers?

---

**Cursor**

Good question. Looking at the text extraction of the output, you can see the watermarks are present -- "CONFIDENTIAL" on page 1, "DRAFT" on page 2, "INTERNAL USE ONLY" on page 4. But they are **not** on a separate PDF layer (Optional Content Group / OCG).

Currently, the watermark is rendered with `OnTop = false`, meaning it's painted **behind** the page content in the same content stream. This means:

- It **cannot** be toggled on/off in a PDF viewer
- It **cannot** be selectively removed without re-processing
- It's baked into the page rendering

**If you want true PDF layers**, that would be a design change. Let me know if you'd like:

1. **Keep as-is** -- watermark is embedded (simpler, tamper-resistant, good for production documents)
2. **Switch to `OnTop = true`** -- watermark renders *above* content (more visible, still no separate layer)
3. **Add OCG layer support** -- watermark goes into an Optional Content Group so viewers like Adobe Acrobat can show/hide it (requires investigating pdfcpu's OCG support)

Which direction interests you?

---

**User**

Switch to OnTop = true, lets use this

---

**Cursor**

All 31 tests pass. Let me regenerate the output with your NCERT PDF so you can see the difference.

Done. One line changed: `OnTop = false` to `OnTop = true` in `internal/stamp/stamp.go`. The watermarks now render **above** the page content, so they'll be clearly visible overlaying the text. Open `output.pdf` to see the difference.

---

**User**

Generate the pdf again

---

**Cursor**

Done. `output.pdf` has been regenerated. Open it to check the watermarks.

