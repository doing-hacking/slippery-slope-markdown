# Slippery-Slope Markdown
Slippery-Slope Markdown is a Go package which attempts to provide a limited
subset of markdown syntax without causing any clutter. This is useful in
situations where document-level formatting (i.e. headers) is undesirable but
limited text formatting (lists, bold font) could improve readability.

## Features

### Current Features
- Making text bold using `**bold**` syntax
- Making unordered lists using `-` syntax

### Possible Future Features
- Ordered lists?
  - Alternatively, `- 1.` as the canon way to make ordered lists
- Tables?
  - Could be slippery

## Usage
Currently, the only available method is
`ParseNoEscapeFromBytes(w io.Writer, input []byte)`,
which takes a []byte input and writes the parsed output to an io.Writer
(such as a `bytes.Buffer`)

## Inspiration
The main inspiration for this package was that I wanted to see lists in Godoc.
Somebody requested the same feature, but was told it would likely be rejected,
it being a slippery slope to Markdown.
That gave me a great idea for a new Go package - Slippery-Slope Markdown!
This will also be bundled in the doing-hacks godoc fork.
