package sem_test

import (
	"log"
	"testing"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/gocc/scanner"
	"github.com/mewmew/uc/sem"
)

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
			want: `79: redefinition of "a" with different type: "char" vs "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se05.c",
			want: `79: redefinition of "a" with different type: "void(void)" vs "int"`,
		},
		{
			path: "../testdata/incorrect/semantic/se06.c",
			want: `104: redefinition of "a"; previously declared at 70`,
		},
		{
			path: "../testdata/incorrect/semantic/se07.c",
			want: `91: returning "int" from a function with incompatible result type "void"`,
		},
		{
			path: "../testdata/incorrect/semantic/se08.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se09.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se10.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se11.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se12.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se13.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se14.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se15.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se16.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se17.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se18.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se19.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se20.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se21.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se22.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se23.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se24.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se25.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se26.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se27.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se28.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se29.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se30.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se31.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se32.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se33.c",
			want: "foo",
		},
		{
			path: "../testdata/incorrect/semantic/se34.c",
			want: "foo",
		},
	}

	for _, g := range golden {
		log.Println("path:", g.path)
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
