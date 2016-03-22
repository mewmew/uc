package lexer

import (
	"sort"

	"github.com/mewmew/uc/uc/token"
)

// keywords is the set of valid keywords in the µC programming language, sorted
// alphabetically.
var keywords []string

func init() {
	keywords = make([]string, 0, len(token.Keywords))
	for keyword := range token.Keywords {
		keywords = append(keywords, keyword)
	}
	sort.Strings(keywords)
}
