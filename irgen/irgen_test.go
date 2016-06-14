package irgen_test

import (
	"fmt"
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
		// Global variable declarations.
		{
			path: "../testdata/extra/irgen/global_def.c",
			want: "../testdata/extra/irgen/global_def.ll",
		},
		{
			path: "../testdata/extra/irgen/tentative_def.c",
			want: "../testdata/extra/irgen/tentative_def.ll",
		},
		// Return statements.
		{
			path: "../testdata/extra/irgen/void_ret.c",
			want: "../testdata/extra/irgen/void_ret.ll",
		},
		{
			path: "../testdata/extra/irgen/implicit_void_ret.c",
			want: "../testdata/extra/irgen/implicit_void_ret.ll",
		},
		{
			path: "../testdata/extra/irgen/int_ret.c",
			want: "../testdata/extra/irgen/int_ret.ll",
		},
		{
			path: "../testdata/extra/irgen/expr_ret.c",
			want: "../testdata/extra/irgen/expr_ret.ll",
		},
		// Local variable declarations.
		{
			path: "../testdata/extra/irgen/local_def.c",
			want: "../testdata/extra/irgen/local_def.ll",
		},
		// If statements.
		{
			path: "../testdata/extra/irgen/if_stmt.c",
			want: "../testdata/extra/irgen/if_stmt.ll",
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
		info, err := sem.Check(file)
		if err != nil {
			t.Errorf("%q: semantic analysis error: %v", g.path, err)
			continue
		}

		// Generate IR.
		// TODO: Remove debug output.
		fmt.Printf("\n=== [ %s ] ====================================\n\n", g.want)
		module := irgen.Gen(file, info)

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
