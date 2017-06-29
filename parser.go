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
// Ordered List
// 1. This is item 1
// 2. This is item 2
// 4. Inconsistent numbering gets ignored
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

// AddUListItem creates an unordered list item, starting a new list if necessary.
func (parser *BytesToWriterParser) AddUListItem() {
	if parser.UListState {
		parser.Writer.Write([]byte("</li><li>"))
	} else {
		parser.Writer.Write([]byte("<ul><li>"))
		parser.UListState = true
	}
}

// AddUListEnd closes an unordered list.
func (parser *BytesToWriterParser) AddUListEnd() {
	parser.Writer.Write([]byte("</li></ul>"))
	parser.UListState = false
}

// AddOListItem creates an ordered list item, starting a new list if necessary.
func (parser *BytesToWriterParser) AddOListItem() {
	if parser.OListState {
		parser.Writer.Write([]byte("</li><li>"))
	} else {
		parser.Writer.Write([]byte("<ol><li>"))
		parser.OListState = true
	}
}

// AddOListEnd closes an ordered list.
func (parser *BytesToWriterParser) AddOListEnd() {
	parser.Writer.Write([]byte("</li></ol>"))
	parser.OListState = false
}

// CheckLineType reports if the line is an ordered list item, if the line is
// an unordered list item, and the number of leading spaces on the line.
// The parameter i is assumed to be the index of the line's first character.
func (parser *BytesToWriterParser) CheckLineType(i int) (
	bool, // This is a blank line
	bool, // This is an ordered list
	bool, // This is an unordered list
	int, // Depth of line start
) {
	// parser.Writer.Write(
	// 	[]byte("`i=" + string([]byte{parser.Peek(i)}) + "`"),
	// )
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
		return false, false, false, leadingSpace
	}
	// If a blank line, return blank=true
	if parser.Input[s] == '\n' || parser.Input[s] == '<' {
		return true, false, false, 0
	}
	// parser.Writer.Write([]byte("`s=" + string([]byte{parser.Input[s]}) + "`"))

	// Lists start with "- "
	if parser.Input[s] == '-' && parser.Peek(s+1) == ' ' {
		return false, false, true, leadingSpace + 2
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
		return false, false, false, leadingSpace
	}

	// Add leading digits to start index
	s = s + leadingDigits
	// If the first character is out of range, return false,
	if s >= len(parser.Input) {
		return false, false, false, leadingSpace
	}

	if parser.Input[s] == '.' {
		return false, true, false, leadingSpace + leadingDigits + 2
	}

	return false, false, false, 0
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
			blank, ol, ul, offset := parser.CheckLineType(i)

			if blank {
				if parser.UListState {
					parser.AddUListEnd()
				}
				if parser.OListState {
					parser.AddOListEnd()
				}
				w.Write([]byte{'\n'})
			} else {
				if ul {
					parser.AddUListItem()
				} else {
					if parser.UListState {
						w.Write([]byte{' '})
					}
				}

				if ol {
					parser.AddOListItem()
				} else {
					if parser.OListState {
						w.Write([]byte{' '})
					}
				}
			}

			if ul || ol {
				i += offset
			}
			w.Write([]byte{input[i]})

			parseState = StateNormal
			break STATE
		}
		i++
	}
}
