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
	// Input source position (in bytes).
	Pos int
	// Error message.
	Text string
	// Input source.
	Src *Source
}

// New returns a new error based on the given positional information (offset in
// bytes).
func New(pos int, text string) *Error {
	err := &Error{
		Pos:  pos,
		Text: text,
	}
	return err
}

// Newf returns a new formatted error based on the given positional information
// (offset in bytes).
func Newf(pos int, format string, a ...interface{}) *Error {
	err := &Error{
		Pos:  pos,
		Text: fmt.Sprintf(format, a...),
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
	pos := fmt.Sprintf("(byte offset %d)", e.Pos)
	prefix := "error:"
	text := e.Text
	if UseColor {
		pos = term.Color(pos, term.Bold)
		prefix = term.RedBold(prefix)
		text = term.Color(text, term.Bold)
	}
	src := e.Src
	if src == nil {
		// If Src is nil, the error format is as follows.
		//
		//    (byte offset %d) error: text
		return fmt.Sprintf("%s %s %s", pos, prefix, text)
	}
	// The error format is as follows.
	//
	//    (file:line) error: text
	//       1 = y
	//         ^
	line, col := src.Position(e.Pos)
	srcLine := src.Input[src.Lines[line-1]:src.Lines[line]]
	srcLine = strings.Replace(srcLine, "\t", " ", -1)
	srcLine = strings.TrimRight(srcLine, "\n\r")
	arrow := fmt.Sprintf("%*s", col, "^")
	pos = fmt.Sprintf("(%s:%d)", src.Path, line)
	if UseColor {
		pos = term.Color(pos, term.Bold)
		arrow = term.Color(arrow, term.Bold)
	}
	return fmt.Sprintf("%s %s %s\n%s\n%s", pos, prefix, text, srcLine, arrow)
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
