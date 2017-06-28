package slipperyslopemd

import "io"

// ParseNoEscapeFromBytes parses a byte slice to Slippery-Slope Markdown without
// escaping any existing HTML characters.
func ParseNoEscapeFromBytes(w io.Writer, input []byte) {
	const (
		// Parsing states
		StateNormal         = 0
		StateAsterisk       = 1
		StateDoubleAsterisk = 2

		// Format states
		FormatStateNormal = 0
		FormatStateBold   = 1
	)

	parseState := 0
	formatState := 0

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
				if b == '*' {
					parseState = StateAsterisk
					break STATE
				}
				w.Write([]byte{b})
			}
		case StateAsterisk:
			b := input[i]
			if b == '*' {
				switch formatState {
				case FormatStateNormal:
					formatState = FormatStateBold
					w.Write([]byte("<b>"))
				case FormatStateBold:
					formatState = FormatStateNormal
					w.Write([]byte("</b>"))
				}
			} else {
				w.Write([]byte{b})
			}
			parseState = StateNormal
			break STATE
		}
		i++
	}
}
