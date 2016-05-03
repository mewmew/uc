package sem_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/gocc/scanner"
	"github.com/mewmew/uc/sem"
	"github.com/mewmew/uc/sem/errors"
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
		{path: "../testdata/extra/semantic/missing-return-main.c"},
		{path: "../testdata/extra/semantic/tentative-var-def.c"},
		{path: "../testdata/extra/semantic/variable-sized-array-arg.c"},
	}

	errors.UseColor = false

	for _, g := range golden {
		buf, err := ioutil.ReadFile(g.path)
		if err != nil {
			t.Errorf("%q: %v", g.path, err)
			continue
		}
		input := string(buf)
		s := scanner.NewFromString(input)
		src := errors.NewSource(g.path, input)

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
				if e, ok := err.(*errors.Error); ok {
					// Unwrap semantic error.
					e.Src = src
				}
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
			want: `(../testdata/incorrect/semantic/se01.c:5) error: undeclared identifier "b"
 a = a + b; // Variable 'b' not defined
         ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se02.c",
			want: `(../testdata/incorrect/semantic/se02.c:5) error: undeclared identifier "foo"
  a = foo(a); // Function 'foo' not defined
      ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se03.c",
			want: `(../testdata/incorrect/semantic/se03.c:3) error: undeclared identifier "output"
  output(0); // Procedure 'output' not defined
  ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se04.c",
			want: `(../testdata/incorrect/semantic/se04.c:5) error: redefinition of "a" with type "char" instead of "int"
char a;  // Redeclaration of 'a'
     ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se05.c",
			want: `(../testdata/incorrect/semantic/se05.c:5) error: redefinition of "a" with type "void(void)" instead of "int"
void a(void) {  // Attempt to redefine variable 'a'
     ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se06.c",
			want: `(../testdata/incorrect/semantic/se06.c:7) error: redefinition of "a"
int a(int i) {   // Redeclaration of 'a'
    ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se07.c",
			want: `(../testdata/incorrect/semantic/se07.c:4) error: returning "int" from a function with incompatible result type "void"
  return 2 * n; // Attempt to return value from procedure
         ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se08.c",
			want: `(../testdata/incorrect/semantic/se08.c:4) error: returning "void" from a function with incompatible result type "int"
  return;  // Void return from function
  ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se09.c",
			want: `(../testdata/incorrect/semantic/se09.c:6) error: returning "char[1]" from a function with incompatible result type "int"
  else return x;    // Return from function with erroneous type
              ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se10.c",
			want: `(../testdata/incorrect/semantic/se10.c:6) error: invalid operation: n[2] (type "int" does not support indexing)
  n[2]; // Index an integer
   ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se11.c",
			want: `(../testdata/incorrect/semantic/se11.c:4) error: cannot assign to "a" of type "int(void)"
  a = 1; // 'a' is not an lval
    ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se12.c",
			want: `(../testdata/incorrect/semantic/se12.c:6) error: cannot call non-function "a" of type "int"
  a(2); // 'a' is not a function
   ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se13.c",
			want: `(../testdata/incorrect/semantic/se13.c:8) error: invalid operands to binary expression: 1 + foo(0) ("int" and "void")
  1 + foo(0); // 'foo' does not return a value
    ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se14.c",
			want: `(../testdata/incorrect/semantic/se14.c:12) error: cannot call non-function "f" of type "int"
  f(n);  // 'f' refers only to the local variable
   ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se15.c",
			want: `(../testdata/incorrect/semantic/se15.c:8) error: calling "q" with too few arguments; expected 3, got 2
  1 + q(1, 3); // Too few arguments to function 'q'
       ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se16.c",
			want: `(../testdata/incorrect/semantic/se16.c:9) error: calling "d" with too many arguments; expected 2, got 3
  d(1, 2, 3); // Too many arguments to function 'd'
   ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se17.c",
			want: `(../testdata/incorrect/semantic/se17.c:6) error: invalid operation: hello + 1 (type mismatch between "char[5]" and "int")
  hello+1; //  Attempt to use char array in arithmetic. (legal in C)
       ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se18.c",
			want: `(../testdata/incorrect/semantic/se18.c:6) error: cannot assign to "a" of type "char[10]"
  a = 42;   // assign int to array of char
    ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se19.c",
			want: `(../testdata/incorrect/semantic/se19.c:5) error: invalid operation: a == 42 (type mismatch between "char[10]" and "int")
  if (a==42) ;
       ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se20.c",
			want: `(../testdata/incorrect/semantic/se20.c:7) error: cannot assign to "a" of type "int[10]"
  a=b;
   ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se21.c",
			want: `(../testdata/incorrect/semantic/se21.c:5) error: returning "char[10]" from a function with incompatible result type "int"
    return bv;  //  Return from function with erroneous type
           ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se22.c",
			want: `(../testdata/incorrect/semantic/se22.c:6) error: invalid operation: a + 1 (type mismatch between "char[10]" and "int")
  a+1; // Attempt to apply arithmetic to array reference
   ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se23.c",
			want: `(../testdata/incorrect/semantic/se23.c:6) error: invalid operation: b[0] (type "int" does not support indexing)
  return b[0]; //not an array!
          ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se24.c",
			want: `(../testdata/incorrect/semantic/se24.c:6) error: cannot assign to "b" of type "int[10]"
  b = a;  // b cannot be assigned
    ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se25.c",
			want: `(../testdata/incorrect/semantic/se25.c:4) error: cannot assign to "(1 + 2)" of type "int"
  (1 + 2) = 3; //No assignment here!
          ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se26.c",
			want: `(../testdata/incorrect/semantic/se26.c:9) error: calling "f" with incompatible argument type "char[10]" to parameter of type "int[]"
  f(a);
    ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se27.c",
			want: `(../testdata/incorrect/semantic/se27.c:4) error: returning "int" from a function with incompatible result type "void"
  if (1<2) return 2 * n; // Attempt to return value from procedure
                  ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se28.c",
			want: `(../testdata/incorrect/semantic/se28.c:5) error: returning "int" from a function with incompatible result type "void"
  else  return 2 * n; // Attempt to return value from procedure
               ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se29.c",
			want: `(../testdata/incorrect/semantic/se29.c:4) error: redefinition of "n" with type "char" instead of "int"
  char n;
       ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se30.c",
			want: `(../testdata/incorrect/semantic/se30.c:6) error: cannot assign to "a" (type mismatch between "int" and "int[10]")
  a=b;
   ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se31.c",
			want: `(../testdata/incorrect/semantic/se31.c:5) error: redefinition of "a" with type "void(void)" instead of "int"
void a(void);   // Attempt to redefine  'a' as extern
     ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se32.c",
			want: `(../testdata/incorrect/semantic/se32.c:6) error: invalid operands to binary expression: 1 + foo(0) ("int" and "void")
  1 + foo(0); // 'foo' does not return a value
    ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se33.c",
			want: `(../testdata/incorrect/semantic/se33.c:6) error: calling "q" with too few arguments; expected 3, got 2
  1 + q(1, 3); // Too few arguments to function 'q'
       ^`,
		},
		{
			path: "../testdata/incorrect/semantic/se34.c",
			want: `(../testdata/incorrect/semantic/se34.c:6) error: calling "d" with too many arguments; expected 2, got 3
  d(1, 2, 3); // Too many arguments to function 'd'
   ^`,
		},

		// Extra test cases.
		{
			path: "../testdata/extra/semantic/extra-void-arg.c",
			want: `(../testdata/extra/semantic/extra-void-arg.c:4) error: "void" must be the only parameter
void f(int a, void) {
      ^`,
		},
		{
			path: "../testdata/extra/semantic/incompatible-arg-type.c",
			want: `(../testdata/extra/semantic/incompatible-arg-type.c:10) error: calling "a" with incompatible argument type "int" to parameter of type "int[]"
 return a(b);
          ^`,
		},
		{
			path: "../testdata/extra/semantic/index-array.c",
			want: `(../testdata/extra/semantic/index-array.c:7) error: invalid array index; expected integer, got "int[20]"
 x[y];
   ^`,
		},
		{
			path: "../testdata/extra/semantic/local-var-redef.c",
			want: `(../testdata/extra/semantic/local-var-redef.c:6) error: redefinition of "x"
 int x;
     ^`,
		},
		{
			path: "../testdata/extra/semantic/missing-return.c",
			want: `(../testdata/extra/semantic/missing-return.c:10) error: missing return at end of non-void function "f"
}
^`,
		},
		{
			path: "../testdata/extra/semantic/param-redef.c",
			want: `(../testdata/extra/semantic/param-redef.c:5) error: redefinition of "x"
 int x;
     ^`,
		},
		{
			path: "../testdata/extra/semantic/unnamed-arg.c",
			want: `(../testdata/extra/semantic/unnamed-arg.c:4) error: parameter name obmitted
void f(int) {
       ^`,
		},
		{
			path: "../testdata/extra/semantic/variable-sized-array.c",
			want: `(../testdata/extra/semantic/variable-sized-array.c:5) error: array size or initializer missing for "y"
 char y[];
      ^`,
		},
		{
			path: "../testdata/extra/semantic/void-arg.c",
			want: `(../testdata/extra/semantic/void-arg.c:4) error: "x" has invalid type "void"
void f(void x) {
            ^`,
		},
		{
			path: "../testdata/extra/semantic/void-array.c",
			want: `(../testdata/extra/semantic/void-array.c:5) error: invalid element type "void" of array "x"
 void x[10];
      ^`,
		},
		{
			path: "../testdata/extra/semantic/void-array-arg.c",
			want: `(../testdata/extra/semantic/void-array-arg.c:4) error: invalid element type "void" of array "x"
void f(void x[]) {
            ^`,
		},
		{
			path: "../testdata/extra/semantic/void-params.c",
			want: `(../testdata/extra/semantic/void-params.c:4) error: "void" must be the only parameter
void f(void, void) {
      ^`,
		},
		{
			path: "../testdata/extra/semantic/void-var.c",
			want: `(../testdata/extra/semantic/void-var.c:5) error: "x" has invalid type "void"
 void x;
      ^`,
		},
	}

	errors.UseColor = false

	for _, g := range golden {
		buf, err := ioutil.ReadFile(g.path)
		if err != nil {
			t.Errorf("%q: %v", g.path, err)
			continue
		}
		input := string(buf)
		s := scanner.NewFromString(input)
		src := errors.NewSource(g.path, input)

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
				if e, ok := err.(*errors.Error); ok {
					// Unwrap semantic error.
					e.Src = src
				}
			}
			got = err.Error()
		}
		if got != g.want {
			t.Errorf("%q: error mismatch; expected `%v`, got `%v`", g.path, g.want, got)
		} else if strings.Contains(g.path, "extra") {
			// TODO: Remove once sem passes the extra tests.
			fmt.Println("PASS:", g.path)
		}
	}
}

// TODO: add benchmark
