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
}

// Peek returns the byte at index i, or the byte 'X' if i is out of range.
func (parser *BytesToWriterParser) Peek(i int) byte {
	if i >= len(parser.Input) {
		return 'X'
	}
	return parser.Input[i]
}

// ParseNoEscapeFromBytes parses a byte slice to Slippery-Slope Markdown without
// escaping any existing HTML characters.
func ParseNoEscapeFromBytes(w io.Writer, input []byte) {
	parser := &BytesToWriterParser{w, input, 0, false, false, false}

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
					parseState = StateAsterisk
					break STATE
				case '\n':
					parseState = StateLineFeed
					break STATE
				}
				w.Write([]byte{b})
			}
		case StateAsterisk:
			b := input[i]
			if b == '*' {
				if parser.BoldState {
					w.Write([]byte("</b>"))
					parser.BoldState = false
				} else {
					w.Write([]byte("<b>"))
					parser.BoldState = true
				}
			} else {
				w.Write([]byte{b})
			}
			parseState = StateNormal
			break STATE
		case StateLineFeed:
			b := input[i]
			if b == '-' {
				if parser.UListState {
					w.Write([]byte("</li><li>"))
				} else {
					w.Write([]byte("<ul><li>"))
					parser.UListState = true
				}
			} else {
				if parser.UListState {
					w.Write([]byte("</li></ul>"))
					parser.UListState = false
				}
				w.Write([]byte{b})
			}
			parseState = StateNormal
			break STATE
		}
		i++
	}
}
