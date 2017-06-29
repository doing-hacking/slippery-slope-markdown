// Package slipperyslopemd provides a slipper slope to Markdown without slipping
//
// Lists
//
// Unordered look like this
// - Item 1
// - Item 2
// which has another line below it
// - Item 3
//
// Not a list
// - Consecutive list
//
// Toggling bold text **on and** off and** on again and then **off again.
package slipperyslopemd

import "io"

// BytesToWriterParser converts a byte slice containing Slippery-Slope Markdown
// formatted text and converts it to HTML.
type BytesToWriterParser struct {
	Writer     io.Writer
	Input      []byte
	Index      int
	BoldState  bool
	UListState bool
	OListState bool
	ListDepth  int
}

// Peek returns the byte at index i, or the byte 'X' if i is out of range.
func (parser *BytesToWriterParser) Peek(i int) byte {
	if i >= len(parser.Input) {
		return 'X'
	}
	return parser.Input[i]
}

// AddToggleBold adds an open or closing <b> tag based on the state.
func (parser *BytesToWriterParser) AddToggleBold() {
	if parser.BoldState {
		parser.Writer.Write([]byte("</b>"))
		parser.BoldState = false
	} else {
		parser.Writer.Write([]byte("<b>"))
		parser.BoldState = true
	}
}

// CheckLineType reports if the line is an ordered list item, if the line is
// an unordered list item, and the number of leading spaces on the line.
// The parameter i is assumed to be the index of the line's first character.
func (parser *BytesToWriterParser) CheckLineType(i int) (
	bool, // This is an ordered list
	bool, // This is an unordered list
	int, // Depth of line start
) {
	// Get the number of leading spaces
	leadingSpace := 0
	for j := i; j < len(parser.Input); j++ {
		if parser.Input[j] == ' ' {
			leadingSpace++
		} else {
			break
		}
	}
	// Get the index of the first character after the leading spaces
	s := i + leadingSpace
	// If the first character is out of range, return
	if s >= len(parser.Input) {
		return false, false, leadingSpace
	}

	// Lists start with "- "
	if parser.Input[s] == '-' && parser.Peek(s+1) == ' ' {
		return false, true, leadingSpace + 2
	}

	// Ordered lists start with any number of consecutive digits plus ". "
	leadingDigits := 0
	for j := s; j < len(parser.Input); j++ {
		if parser.Input[j] >= 0x30 && parser.Input[i] < 0x40 {
			leadingDigits++
		} else {
			break
		}
	}

	if leadingDigits == 0 {
		return false, false, leadingSpace
	}

	// Add leading digits to start index
	s = s + leadingDigits
	// If the first character is out of range, return
	if s >= len(parser.Input) {
		return false, false, leadingSpace
	}

	if parser.Input[s] == '.' {
		return true, false, leadingSpace + leadingDigits + 2
	}

	return false, false, 0
}

// ParseNoEscapeFromBytes parses a byte slice to Slippery-Slope Markdown without
// escaping any existing HTML characters.
func ParseNoEscapeFromBytes(w io.Writer, input []byte) {
	parser := &BytesToWriterParser{w, input, 0, false, false, false, 0}

	const (
		// Parsing states
		StateNormal   = 0
		StateAsterisk = 1
		StateLineFeed = 3
	)

	parseState := 0

	i := 0

	for {
		if !(i < len(input)) {
			break
		}

	STATE:
		switch parseState {
		case StateNormal: // Normal Mode
			for ; i < len(input); i++ {
				b := input[i]
				switch b {
				case '*':
					if parser.Peek(i+1) == '*' {
						parser.AddToggleBold()
						i++
						break STATE
					}
				case '\n':
					parseState = StateLineFeed
					break STATE
				}
				w.Write([]byte{b})
			}
		case StateLineFeed:
			b := input[i]
			_, ul, _ := parser.CheckLineType(i)

			if ul {
				if parser.UListState {
					w.Write([]byte("</li><li>"))
				} else {
					w.Write([]byte("<ul><li>"))
					parser.UListState = true
				}
			} else {
				if parser.UListState {
					if b == '\n' {
						w.Write([]byte("</li></ul>"))
						parser.UListState = false
					} else {
						w.Write([]byte{' '})
					}
				}
				w.Write([]byte{b})
			}
			parseState = StateNormal
			break STATE
		}
		i++
	}
}
