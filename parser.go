package slipperyslopemd

import "io"

// ParseNoEscapeFromBytes parses a byte slice to Slippery-Slope Markdown without
// escaping any existing HTML characters.
func ParseNoEscapeFromBytes(w io.Writer, input []byte) {
	const (
		// Parsing states
		StateNormal   = 0
		StateAsterisk = 1
		StateLineFeed = 3

		// Format states
		FormatStateNormal = 0
		FormatStateBold   = 1

		// List states
		ListStateNotInList = 0
		ListStateInList    = 1
	)

	parseState := 0
	formatState := 0
	listState := 0

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
		case StateLineFeed:
			b := input[i]
			if b == '-' {
				switch listState {
				case ListStateNotInList:
					listState = ListStateInList
					w.Write([]byte("<ul><li>"))
				case ListStateInList:
					w.Write([]byte("</li><li>"))
				}
			} else {
				if listState == ListStateInList {
					w.Write([]byte("</li></ul>"))
				}
				w.Write([]byte{b})
			}
			parseState = StateNormal
			break STATE
		}
		i++
	}
}
