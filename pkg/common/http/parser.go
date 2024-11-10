package common

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/dslipak/pdf"
)

var supportedTextExtensions = map[string]bool{
	".go": true, ".js": true, ".txt": true, ".py": true, ".java": true,
	".cpp": true, ".c": true, ".cs": true, ".rb": true, ".php": true,
	".html": true, ".css": true, ".json": true, ".xml": true, ".sh": true,
	".yml": true, ".md": true, ".ini": true, ".conf": true, ".yaml": true,
	".cfg": true, ".log": true, ".sql": true, ".tsv": true, ".csv": true,
	".bat": true, ".toml": true, ".r": true, ".pl": true, ".m": true,
	".scala": true, ".swift": true, ".vb": true, ".rs": true, ".erl": true,
	".hs": true, ".lhs": true, ".tex": true, ".cls": true, ".sty": true,
	".scss": true, ".sass": true, ".less": true, ".asciidoc": true, ".pdf": true,
	".rst": true, ".org": true, ".srt": true, ".vtt": true, ".xslt": true,
}

func MultipleFileParserToText(files []*multipart.FileHeader) (map[string]string, error) {
	fileMap := make(map[string]string)
	for _, fileHeader := range files {
		content, err := parseFile(fileHeader)
		if err != nil {
			return nil, err
		}
		fileMap[fileHeader.Filename] = content
	}
	return fileMap, nil
}

func parseFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	switch {
	case ext == ".pdf":
		return parsePDFFile(file)
	case supportedTextExtensions[ext]:
		return parseTextFile(file)
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileHeader.Filename)
	}
}

func parseTextFile(file multipart.File) (string, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// parsePDFFile reads and parses a PDF file, returning the extracted text as a string.
func parsePDFFile(file multipart.File) (string, error) {
	// Reset the cursor to the start of the file
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to seek file start: %v", err)
	}

	// Read the file content into memory
	pdfData := &bytes.Buffer{}
	if _, err := io.Copy(pdfData, file); err != nil {
		return "", fmt.Errorf("failed to read PDF file: %v", err)
	}

	// Create a new PDF reader from the byte data
	reader, err := pdf.NewReader(bytes.NewReader(pdfData.Bytes()), int64(pdfData.Len()))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %v", err)
	}

	// Extract text from all pages
	textContent, err := getAllPagesText(reader)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from PDF: %v", err)
	}

	return textContent, nil
}

// Helper function to iterate and extract text from all pages in the PDF.
func getAllPagesText(reader *pdf.Reader) (string, error) {
	var textContent string
	for i := 1; i <= reader.NumPage(); i++ {
		page := reader.Page(i)

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %v", i, err)
		}

		textContent += text + "\n"
	}
	return textContent, nil
}
