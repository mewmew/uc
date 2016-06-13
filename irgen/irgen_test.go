package irgen_test

import (
	"io/ioutil"
	"testing"

	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/gocc/scanner"
	"github.com/mewmew/uc/irgen"
	"github.com/mewmew/uc/sem"
)

func TestGen(t *testing.T) {
	golden := []struct {
		path string
		want string
	}{
		{
			path: "../testdata/extra/irgen/global_def.c",
			want: "../testdata/extra/irgen/global_def.ll",
		},
		{
			path: "../testdata/extra/irgen/tentative_def.c",
			want: "../testdata/extra/irgen/tentative_def.ll",
		},
	}

	for _, g := range golden {
		// Lex input.
		buf, err := ioutil.ReadFile(g.path)
		if err != nil {
			t.Errorf("%q: %v", g.path, err)
			continue
		}
		input := string(buf)
		s := scanner.NewFromString(input)

		// Parse input.
		p := parser.NewParser()
		f, err := p.Parse(s)
		if err != nil {
			t.Errorf("%q: parse error: %v", g.path, err)
			continue
		}
		file := f.(*ast.File)

		// Verify input.
		if err := sem.Check(file); err != nil {
			t.Errorf("%q: semantic analysis error: %v", g.path, err)
			continue
		}

		// Generate IR.
		module, err := irgen.Gen(file)
		if err != nil {
			t.Errorf("%q: unable to generate IR: %v", g.path, err)
			continue
		}

		// Compare generated module against gold standard.
		buf, err = ioutil.ReadFile(g.want)
		if err != nil {
			t.Errorf("%q: %v", g.path, err)
			continue
		}
		got, want := module.String(), string(buf)
		if got != want {
			t.Errorf("%q: module mismatch; expected `%v`, got `%v`", g.path, want, got)
		}
	}
}
