package pdf

import (
	"bytes"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

// Parser handles PDF parsing operations
type Parser struct{}

// NewParser creates a new PDF parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseResult represents the result of parsing a PDF
type ParseResult struct {
	Text       string
	PageCount  int
	Metadata   map[string]string
	Title      string
	Authors    []string
	Abstract   string
	References []string
}

// Parse extracts text and metadata from a PDF
func (p *Parser) Parse(reader io.ReaderAt, size int64) (*ParseResult, error) {
	pdfReader, err := pdf.NewReader(reader, size)
	if err != nil {
		return nil, err
	}

	result := &ParseResult{
		Metadata: make(map[string]string),
		PageCount: pdfReader.NumPage(),
	}

	// Extract text from all pages
	var textBuilder strings.Builder
	for i := 1; i <= pdfReader.NumPage(); i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		textBuilder.WriteString(text)
		textBuilder.WriteString("\n\n")
	}

	result.Text = textBuilder.String()

	// Try to extract metadata
	result.Title = extractTitle(result.Text)
	result.Authors = extractAuthors(result.Text)
	result.Abstract = extractAbstract(result.Text)
	result.References = extractReferences(result.Text)

	return result, nil
}

// ParseFromBytes parses a PDF from a byte slice
func (p *Parser) ParseFromBytes(data []byte) (*ParseResult, error) {
	reader := bytes.NewReader(data)
	return p.Parse(reader, int64(len(data)))
}

// extractTitle attempts to extract the title from PDF text
func extractTitle(text string) string {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and very short lines
		if len(line) > 10 && len(line) < 200 {
			// Title is usually the first substantial line
			// that doesn't look like an author name or institution
			if !containsCommonNonTitlePatterns(line) {
				return line
			}
		}
	}
	return ""
}

// extractAuthors attempts to extract authors from PDF text
func extractAuthors(text string) []string {
	var authors []string
	lines := strings.Split(text, "\n")
	
	foundTitle := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 5 {
			continue
		}

		// Skip until we find a title-like line
		if !foundTitle {
			if len(line) > 10 && !containsCommonNonTitlePatterns(line) {
				foundTitle = true
			}
			continue
		}

		// After title, look for author patterns
		// Authors are often comma-separated or on separate lines
		if containsAuthorPatterns(line) {
			// Split by common separators
			parts := strings.FieldsFunc(line, func(r rune) bool {
				return r == ',' || r == ';' || r == '·'
			})
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if looksLikeName(part) {
					authors = append(authors, part)
				}
			}
			// Only check a few lines after title
			if len(authors) > 0 {
				break
			}
		}
		
		// Stop after checking some lines
		if len(authors) == 0 && foundTitle {
			break
		}
	}

	return authors
}

// extractAbstract attempts to extract the abstract from PDF text
func extractAbstract(text string) string {
	// Look for "Abstract" keyword
	lowerText := strings.ToLower(text)
	abstractIdx := strings.Index(lowerText, "abstract")
	if abstractIdx == -1 {
		return ""
	}

	// Find the start of abstract text
	start := abstractIdx + len("abstract")
	for start < len(text) && (text[start] == ' ' || text[start] == '\n' || text[start] == ':' || text[start] == '—') {
		start++
	}

	// Find the end (usually "Introduction", "Keywords", or double newline)
	endMarkers := []string{"introduction", "1.", "keywords", "key words", "1 "}
	end := len(text)
	for _, marker := range endMarkers {
		idx := strings.Index(lowerText[start:], marker)
		if idx != -1 && start+idx < end {
			end = start + idx
		}
	}

	// Limit abstract length
	if end-start > 2000 {
		end = start + 2000
	}

	abstract := strings.TrimSpace(text[start:end])
	
	// Clean up the abstract
	abstract = strings.ReplaceAll(abstract, "\n", " ")
	abstract = strings.Join(strings.Fields(abstract), " ")

	return abstract
}

// extractReferences attempts to extract references from PDF text
func extractReferences(text string) []string {
	var references []string
	
	// Find references section
	lowerText := strings.ToLower(text)
	refIdx := strings.LastIndex(lowerText, "references")
	if refIdx == -1 {
		refIdx = strings.LastIndex(lowerText, "bibliography")
	}
	if refIdx == -1 {
		return references
	}

	refSection := text[refIdx:]
	lines := strings.Split(refSection, "\n")
	
	var currentRef strings.Builder
	for _, line := range lines[1:] { // Skip "References" line
		line = strings.TrimSpace(line)
		if line == "" {
			if currentRef.Len() > 0 {
				references = append(references, currentRef.String())
				currentRef.Reset()
			}
			continue
		}

		// Check if this starts a new reference (usually starts with [n] or n.)
		if startsWithRefNumber(line) && currentRef.Len() > 0 {
			references = append(references, currentRef.String())
			currentRef.Reset()
		}

		if currentRef.Len() > 0 {
			currentRef.WriteString(" ")
		}
		currentRef.WriteString(line)

		// Limit number of references
		if len(references) >= 100 {
			break
		}
	}

	if currentRef.Len() > 0 {
		references = append(references, currentRef.String())
	}

	return references
}

// Helper functions

func containsCommonNonTitlePatterns(line string) bool {
	lowerLine := strings.ToLower(line)
	patterns := []string{
		"university", "institute", "department", "school",
		"@", "email", ".edu", ".com", ".org",
		"abstract", "introduction", "keywords",
	}
	for _, pattern := range patterns {
		if strings.Contains(lowerLine, pattern) {
			return true
		}
	}
	return false
}

func containsAuthorPatterns(line string) bool {
	// Authors often have superscript indicators or email patterns nearby
	// But main indicator is position after title
	return len(line) > 5 && len(line) < 500
}

func looksLikeName(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 2 || len(s) > 50 {
		return false
	}
	
	// Name should have mostly letters and spaces
	letters := 0
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == ' ' || r == '-' || r == '.' {
			letters++
		}
	}
	return float64(letters)/float64(len(s)) > 0.8
}

func startsWithRefNumber(line string) bool {
	// Check for patterns like [1], [2], 1., 2., etc.
	if len(line) < 2 {
		return false
	}
	
	// [n] pattern
	if line[0] == '[' {
		endBracket := strings.Index(line, "]")
		if endBracket > 0 && endBracket < 5 {
			return true
		}
	}
	
	// n. pattern
	if line[0] >= '1' && line[0] <= '9' {
		for i := 1; i < len(line) && i < 4; i++ {
			if line[i] == '.' || line[i] == ')' {
				return true
			}
			if line[i] < '0' || line[i] > '9' {
				break
			}
		}
	}
	
	return false
}

