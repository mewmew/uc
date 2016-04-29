// Package errors provides pretty-printing of semantic analysis errors.
package errors

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mewkiz/pkg/term"
)

// UseColor indicates if error messages should use colors.
var UseColor = true

// An Error represents a semantic analysis error.
type Error struct {
	// Input source position.
	Position int
	// Error message.
	Text string
	// Input source.
	Src *Source
}

// New returns a new error with the given positional information.
func New(pos int, text string) *Error {
	err := &Error{
		Position: pos,
		Text:     text,
	}
	return err
}

// Newf returns a new formatted error with the given positional information.
func Newf(pos int, format string, a ...interface{}) *Error {
	err := &Error{
		Position: pos,
		Text:     fmt.Sprintf(format, a...),
	}
	return err
}

// Error returns an error string with position information.
//
// The error format is as follows.
//
//    (file:line:column): error: text
func (e *Error) Error() string {
	// Use colors.
	pos := e.Position
	sPos := string(pos)
	prefix := "error:"
	text := e.Text
	src := e.Src
	if UseColor {
		sPos = term.Color(sPos, term.Bold)
		prefix = term.RedBold(prefix)
		text = term.Color(text, term.Bold)
	}
	if src == nil {
		// If Src is nil, the format is as follows.
		//
		//    (byte offset %d)
		return fmt.Sprintf("(byte offset %d) %s %s", sPos, prefix, text)
	}
	// The position format is as follows.
	//
	//    (file:line) message type:
	line, col := src.Position(pos)
	srcLine := strings.Replace(src.Input[src.Lines[line-1]:src.Lines[line]], "\t", " ", -1)
	srcLine = strings.Trim(srcLine, "\n\r")
	point := fmt.Sprintf("%*s", col, "^")
	if UseColor {
		point = term.Color(point, term.Bold)
	}
	return fmt.Sprintf("(%s:%d) %s %s\n%s\n%s", src.Path, line, prefix, text, srcLine, point)
}

// A Source represents an input source.
type Source struct {
	// Input source path (file path or <stdin>).
	Path string
	// Input source text.
	Input string
	// Lines tracks the start positions of lines within the input stream.
	Lines []int
}

// Position returns the corresponding line:column pair of the given position in
// the input stream.
func (src *Source) Position(pos int) (line, column int) {
	// Implemented using binary search, lines should be sorted in ascending
	// order.
	index := sort.SearchInts(src.Lines, pos+1) - 1
	line = index + 1
	column = pos - src.Lines[index] + 1
	return line, column
}

// NewSource returns a new source based on the given input. The path is only
// used in error messages, and "<stdin>" is conventionally used for the standard
// input stream.
func NewSource(path, input string) *Source {
	src := &Source{
		Path:  path,
		Input: input,
	}
	for i := 0; i < len(input); {
		src.Lines = append(src.Lines, i)
		pos := strings.IndexRune(input[i:], '\n')
		if pos != -1 {
			i += pos + 1
		}
	}
	return src
}
