// Package pdfmark applies text watermarks to PDF documents.
//
// It reads a PDF template from an io.Reader, watermark instructions from a CSV
// io.Reader, and writes the watermarked PDF to an io.WriteCloser. The function
// is stateless and safe for concurrent use.
//
// CSV format:
//
//	page,watermark_text
//	1,CONFIDENTIAL
//	3,DRAFT
//
// Usage:
//
//	err := pdfmark.Watermark(dst, pdfReader, csvReader)
package pdfmark
