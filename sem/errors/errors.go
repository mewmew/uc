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
	Position
	// Error message.
	Text string
}

// New returns a new error with the given positional information.
func New(pos int, text string) *Error {
	err := &Error{
		Position: Position{
			Pos: pos,
		},
		Text: text,
	}
	return err
}

// Newf returns a new formatted error with the given positional information.
func Newf(pos int, format string, a ...interface{}) *Error {
	err := &Error{
		Position: Position{
			Pos: pos,
		},
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
	pos := e.Position.String()
	prefix := "error:"
	text := e.Text
	if UseColor {
		pos = term.WhiteBold(pos)
		prefix = term.RedBold(prefix)
		text = term.WhiteBold(text)
	}
	return fmt.Sprintf("%s %s %s", pos, prefix, text)
}

// A Position represents an input source position.
type Position struct {
	// Input source position (in bytes).
	Pos int
	// Input source.
	Src *Source
}

// String returns a string representation of the position.
//
// The position format is as follows.
//
//    (file:line:column)
//
// If Src is nil, the position format is as follows.
//
//    (byte offset %d)
func (pos Position) String() string {
	if pos.Src == nil {
		return fmt.Sprintf("(byte offset %d)", pos.Pos)
	}
	line, col := pos.Src.Position(pos.Pos)
	return fmt.Sprintf("(%s:%d:%d)", pos.Src.Path, line, col)
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
