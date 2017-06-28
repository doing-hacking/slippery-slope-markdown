package slipperyslopemd

import (
	"bytes"
	"io"
)

// ParseNoEscapeFromBytes parses a byte slice to Slippery-Slope Markdown without
// escaping any existing HTML characters.
func ParseNoEscapeFromBytes(w io.Writer, input []byte) {
	const lenBuffer = 3
	buffer := []byte{0, 0, 0}
	state := 0

	i := 0

	for {
		if !(i < len(input)) {
			break
		}

	STATE:
		switch state {
		case 0: // Normal Mode
			for ; i < len(input); i++ {
				b := input[i]
				buffer = append(buffer, b)[1:]
				w.Write([]byte{b})
				if bytes.Compare(buffer[lenBuffer-2:], []byte("**")) == 0 {
					w.Write([]byte("<b>"))
					state = 1
					break STATE
				}
			}
		case 1: // Bold Mode
			for ; i < len(input); i++ {
				b := input[i]
				buffer = append(buffer, b)[1:]
				w.Write([]byte{b})
				if bytes.Compare(buffer[lenBuffer-2:], []byte("**")) == 0 {
					w.Write([]byte("</b>"))
					state = 0
					break STATE
				}
			}
		}
		i++
	}
}
