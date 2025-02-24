package parser

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

type Parser struct {
	scanner scanner.Scanner
}

func ConvertToTS(input io.Reader, output io.Writer) error {
	p := &Parser{}
	p.scanner.Init(input)
	p.scanner.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats |
		scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments
	p.scanner.Whitespace = 0

	return p.parse(output)
}

func (p *Parser) writeToken(output io.Writer, text string) {
	_, err := fmt.Fprint(output, text)
	if err != nil {
		fmt.Println("Error writing token:", err)
	}
}

func (p *Parser) parse(output io.Writer) error {
	var isFunctionDeclaration bool
	var isVariableDeclaration bool
	var isParsingFunctionArgs bool
	var isConstructorDefinition bool

	for tok := p.scanner.Scan(); tok != scanner.EOF; tok = p.scanner.Scan() {
		text := p.scanner.TokenText()

		if text == "function" {
			p.writeToken(output, text)
			isFunctionDeclaration = true
			continue
		}

		if text == "constructor" {
			isConstructorDefinition = true
			p.writeToken(output, text)
			continue
		}

		if isConstructorDefinition && text == ")" {
			p.writeToken(output, text)
			isParsingFunctionArgs = false
			continue
		}

		if isConstructorDefinition && text == "}" {
			p.writeToken(output, text)
			isConstructorDefinition = false
			continue
		}

		if text == "var" || text == "const" || text == "let" {
			p.writeToken(output, text)
			isVariableDeclaration = true
			continue
		}

		if isVariableDeclaration && text == "=" {
			p.writeToken(output, ": any ")
			p.writeToken(output, text)
			isVariableDeclaration = false
			continue
		}

		if isFunctionDeclaration && text == "}" {
			p.writeToken(output, text)
			isFunctionDeclaration = false
			continue
		}

		if text == "(" {
			p.writeToken(output, text)
			isParsingFunctionArgs = true
			continue
		}

		if isParsingFunctionArgs && text == ")" {
			p.writeToken(output, text)
			p.writeToken(output, ": any")
			isParsingFunctionArgs = false
			continue
		}

		if isParsingFunctionArgs &&
			text != "," &&
			strings.TrimSpace(text) != "" {
			p.writeToken(output, text)
			p.writeToken(output, ": any")
			continue
		}

		p.writeToken(output, text)
	}

	return nil
}
