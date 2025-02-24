package parser

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"
	"unicode"
)

type Parser struct {
	scanner     scanner.Scanner
	buffer      strings.Builder
	lastRune    rune
	col         int    // current column
	line        int    // current line
	lastNewline bool   // track if we just processed a newline
	indent      string // current indentation
}

func ConvertToJS(input io.Reader, output io.Writer) error {
	p := &Parser{}
	p.scanner.Init(input)
	p.scanner.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats |
		scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments
	p.scanner.Whitespace = 0 // Don't skip any whitespace
	p.line = 1

	return p.parse(output)
}

func (p *Parser) writeToken(output io.Writer, text string) {
	// Handle whitespace preservation
	pos := p.scanner.Position
	if pos.Line > p.line {
		// Preserve newlines
		newlines := strings.Repeat("\n", pos.Line-p.line)
		fmt.Fprint(output, newlines)
		p.line = pos.Line
		p.col = 1
		p.lastNewline = true
	}

	// Preserve horizontal spacing
	if p.lastNewline {
		// Preserve indentation after newline
		if len(text) > 0 && text[0] != '\n' {
			spaces := strings.Repeat(" ", pos.Column-1)
			fmt.Fprint(output, spaces)
		}
		p.lastNewline = false
	} else {
		// Preserve spacing between tokens
		if pos.Column > p.col {
			spaces := strings.Repeat(" ", pos.Column-p.col)
			fmt.Fprint(output, spaces)
		}
	}

	fmt.Fprint(output, text)
	p.col = pos.Column + len(text)
}

func (p *Parser) parse(output io.Writer) error {
	var inTypeDecl, inInterface bool
	var bracketCount int
	var skipUntilBracket bool

	for tok := p.scanner.Scan(); tok != scanner.EOF; tok = p.scanner.Scan() {
		text := p.scanner.TokenText()

		// Handle import statements
		if text == "import" {
			p.writeToken(output, text)
			// Copy the import statement as-is until semicolon
			for tok := p.scanner.Scan(); tok != scanner.EOF; tok = p.scanner.Scan() {
				text = p.scanner.TokenText()
				p.writeToken(output, text)
				if text == ";" {
					break
				}
			}
			continue
		}

		// Handle keywords that might start type-related declarations
		switch text {
		case "interface":
			inInterface = true
			skipUntilBracket = true
			continue
		case "type":
			inTypeDecl = true
			skipUntilBracket = true
			continue
		}

		// Track brackets for block scope
		if text == "{" {
			bracketCount++
			if skipUntilBracket {
				continue
			}
		} else if text == "}" {
			bracketCount--
			if bracketCount == 0 {
				inInterface = false
				inTypeDecl = false
				skipUntilBracket = false
				continue
			}
		}

		// Skip content if we're in a type declaration or interface
		if inInterface || inTypeDecl || skipUntilBracket {
			continue
		}

		// Handle type annotations
		if text == ":" && !isInString(p.lastRune) {
			// Skip until we find a semicolon, equals sign, or opening brace
			for tok := p.scanner.Scan(); tok != scanner.EOF; tok = p.scanner.Scan() {
				text = p.scanner.TokenText()
				if text == ";" || text == "=" || text == "{" {
					if text == "{" {
						bracketCount++
					}
					p.writeToken(output, text)
					break
				}
			}
			continue
		}

		// Handle generics
		if text == "<" {
			if isTypeParameter(p.scanner.TokenText()) {
				// Skip until matching ">"
				bracketCount = 1
				for tok := p.scanner.Scan(); tok != scanner.EOF; tok = p.scanner.Scan() {
					if p.scanner.TokenText() == ">" {
						bracketCount--
						if bracketCount == 0 {
							break
						}
					} else if p.scanner.TokenText() == "<" {
						bracketCount++
					}
				}
				continue
			} else {
				// Not a generic type parameter, write both tokens
				p.writeToken(output, text)
				p.writeToken(output, p.scanner.TokenText())
			}
			continue
		}

		// Write the token
		p.writeToken(output, text)
		if len(text) > 0 {
			p.lastRune = rune(text[len(text)-1])
		}
	}

	return nil
}

func isInString(r rune) bool {
	return r == '"' || r == '\''
}

func isTypeParameter(s string) bool {
	return unicode.IsUpper(rune(s[0])) || s == "extends" || s == "keyof"
}
