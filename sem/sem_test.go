package sem_test

import (
	"testing"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/gocc/scanner"
	"github.com/mewmew/uc/sem"
)

func TestCheckValid(t *testing.T) {
	var golden = []struct {
		path string
	}{
		{path: "../testdata/quiet/semantic/s01.c"},
		{path: "../testdata/quiet/semantic/s02.c"},
		{path: "../testdata/quiet/semantic/s03.c"},
		{path: "../testdata/quiet/semantic/s04.c"},
		{path: "../testdata/quiet/semantic/s05.c"},
		{path: "../testdata/quiet/semantic/s06.c"},
	}

	for _, g := range golden {
		s, err := scanner.Open(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		p := parser.NewParser()
		file, err := p.Parse(s)
		if err != nil {
			t.Error(err)
			continue
		}
		f := file.(*ast.File)
		err = sem.Check(f)
		if err != nil {
			if e, ok := err.(*errutil.ErrInfo); ok {
				// Unwrap errutil error.
				err = e.Err
			}
			t.Errorf("%q: unexpected error: `%v`", g.path, err.Error())
		}
	}
}

func TestCheckError(t *testing.T) {
	var golden = []struct {
		path string
		want string
	}{
		{
			path: "../testdata/incorrect/semantic/se01.c",
			want: `100: undeclared identifier "b"`,
		},
		{
			path: "../testdata/incorrect/semantic/se02.c",
			want: `96: undeclared identifier "foo"`,
		},
		{
			path: "../testdata/incorrect/semantic/se03.c",
			want: `84: undeclared identifier "output"`,
		},
		{
			path: "../testdata/incorrect/semantic/se04.c",
			want: `79: redefinition of "a" with type "char" instead of "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se05.c",
			want: `79: redefinition of "a" with type "void(void)" instead of "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se06.c",
			want: `104: redefinition of "a"; previously defined at 70`,
		},
		{
			path: "../testdata/incorrect/semantic/se07.c",
			want: `91: returning "int" from a function with incompatible result type "void"`,
		},
		{
			path: "../testdata/incorrect/semantic/se08.c",
			want: `83: returning "void" from a function with incompatible result type "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se09.c",
			want: `132: returning "char[1]" from a function with incompatible result type "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se10.c",
			want: `103: invalid operation: n[2] (type "int" does not support indexing)`,
		},
		{
			path: "../testdata/incorrect/semantic/se11.c",
			want: `84: cannot assign to "a" of type "int(void)"`,
		},
		{
			path: "../testdata/incorrect/semantic/se12.c",
			want: `94: cannot call non-function "a" of type "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se13.c",
			want: `112: invalid operands to binary expression: 1 + foo(0) ("int" and "void")`,
		},
		{
			path: "../testdata/incorrect/semantic/se14.c",
			want: `143: cannot call non-function "f" of type "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se15.c",
			want: `147: calling "q" with too few arguments; expected 3, got 2`,
		},
		{
			path: "../testdata/incorrect/semantic/se16.c",
			want: `128: calling "d" with too many arguments; expected 2, got 3`,
		},
		{
			path: "../testdata/incorrect/semantic/se17.c",
			want: `108: invalid operation: hello + 1 (type mismatch between "char[5]" and "int")`,
		},
		{
			path: "../testdata/incorrect/semantic/se18.c",
			want: `101: cannot assign to "a" of type "char[10]"`,
		},
		{
			path: "../testdata/incorrect/semantic/se19.c",
			want: `103: invalid operation: a == 42 (type mismatch between "char[10]" and "int")`,
		},
		{
			path: "../testdata/incorrect/semantic/se20.c",
			want: `113: cannot assign to "a" of type "int[10]"`,
		},
		{
			path: "../testdata/incorrect/semantic/se21.c",
			want: `107: returning "char[10]" from a function with incompatible result type "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se22.c",
			want: `100: invalid operation: a + 1 (type mismatch between "char[10]" and "int")`,
		},
		{
			path: "../testdata/incorrect/semantic/se23.c",
			want: `97: invalid operation: b[0] (type "int" does not support indexing)`,
		},
		{
			path: "../testdata/incorrect/semantic/se24.c",
			want: `106: cannot assign to "b" of type "int[10]"`,
		},
		{
			path: "../testdata/incorrect/semantic/se25.c",
			want: `94: cannot assign to "(1 + 2)" of type "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se26.c",
			want: `132: calling "f" with incompatible argument type "char[10]" to parameter of type "int[]"`,
		},
		{
			path: "../testdata/incorrect/semantic/se27.c",
			want: `101: returning "int" from a function with incompatible result type "void"`,
		},
		{
			path: "../testdata/incorrect/semantic/se28.c",
			want: `113: returning "int" from a function with incompatible result type "void"`,
		},
		{
			path: "../testdata/incorrect/semantic/se29.c",
			want: `90: redefinition of "n" with type "char" instead of "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se30.c",
			want: `47: cannot assign to "a" (type mismatch between "int" and "int[10]")`,
		},
		{
			path: "../testdata/incorrect/semantic/se31.c",
			want: `79: redefinition of "a" with type "void(void)" instead of "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se32.c",
			want: `105: invalid operands to binary expression: 1 + foo(0) ("int" and "void")`,
		},
		{
			path: "../testdata/incorrect/semantic/se33.c",
			want: `119: calling "q" with too few arguments; expected 3, got 2`,
		},
		{
			path: "../testdata/incorrect/semantic/se34.c",
			want: `109: calling "d" with too many arguments; expected 2, got 3`,
		},
	}

	for _, g := range golden {
		s, err := scanner.Open(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		p := parser.NewParser()
		file, err := p.Parse(s)
		if err != nil {
			t.Error(err)
			continue
		}
		f := file.(*ast.File)

		err = sem.Check(f)
		got := ""
		if err != nil {
			if e, ok := err.(*errutil.ErrInfo); ok {
				// Unwrap errutil error.
				err = e.Err
			}
			got = err.Error()
		}
		if got != g.want {
			t.Errorf("%q: error mismatch; expected `%v`, got `%v`", g.path, g.want, got)
		}
	}
}

// TODO: add benchmark
