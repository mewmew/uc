// Package ast declares the types used to represent abstract syntax trees of µC
// soure code.
package ast

import (
	"bytes"
	"fmt"

	"github.com/mewmew/uc/token"
	"github.com/mewmew/uc/types"
)

// A File represents a µC source file.
//
// Examples.
//
//    int x; int main(void) { x = 42; return x; }
type File struct {
	// Top-level declarations.
	Decls []Decl
}

// A Node represents a node within the abstract syntax tree, and has one of the
// following underlying types.
//
//    *File
//    Decl
//    Stmt
//    Expr
//    Type
type Node interface {
	fmt.Stringer
	// Start returns the start position of the node within the input stream.
	Start() int
}

// A Decl node represents a declaration, and has one of the following underlying
// types.
//
//    *FuncDecl
//    *VarDecl
//    *TypeDef
//
// Pseudo-code representation of a declaration.
//
//    type ident [= value]
type Decl interface {
	Node
	// Type returns the type of the declared identifier.
	Type() types.Type
	// Name returns the name of the declared identifier.
	Name() *Ident
	// Value returns the initializing value of the defined identifier; or nil if
	// declaration or tentative definition.
	//
	// Underlying type for function declarations.
	//
	//    *BlockStmt
	//
	// Underlying type for variable declarations.
	//
	//    Expr
	//
	// Underlying type for type definitions.
	//
	//    Type
	Value() Node
	// isDecl ensures that only declaration nodes can be assigned to the Decl
	// interface.
	isDecl()
}

// Declaration nodes.
type (
	// A FuncDecl node represents a function declaration.
	//
	// Examples.
	//
	//    int puts(char s[]);
	//    int add(int a, int b) { return a+b; }
	FuncDecl struct {
		// Function signature.
		FuncType *FuncType
		// Function name.
		FuncName *Ident
		// Function body; or nil if function declaration (i.e. not function
		// definition).
		Body *BlockStmt
	}

	// A VarDecl node represents a variable declaration.
	//
	// Examples.
	//
	//    int x;
	//    char buf[128];
	VarDecl struct {
		// Variable type.
		VarType Type
		// Variable name.
		VarName *Ident
		// Variable value expression; or nil if variable declaration (i.e. not
		// variable definition).
		Val Expr
	}

	// A TypeDef node represents a type definition.
	//
	// Examples.
	//
	//    typedef int foo;
	TypeDef struct {
		// Position of `typedef` keyword.
		Typedef int
		// Underlying type of type definition.
		DeclType Type
		// Type name.
		TypeName *Ident
	}
)

// A Stmt node represents a statement, and has one of the following underlying
// types.
//
//    *BlockStmt
//    *EmptyStmt
//    *ExprStmt
//    *IfStmt
//    *ReturnStmt
//    *WhileStmt
type Stmt interface {
	Node
	// isStmt ensures that only statement nodes can be assigned to the Stmt
	// interface.
	isStmt()
}

// Statement nodes.
type (
	// A BlockStmt node represents a block statement.
	//
	// Examples.
	//
	//    {}
	//    { int x; x = 42; }
	BlockStmt struct {
		// Position of left-brace `{`.
		Lbrace int
		// List of block items contained within the block.
		Items []BlockItem
		// Position of right-brace `}`.
		Rbrace int
	}

	// An EmptyStmt node represents an empty statement (i.e. ";").
	//
	// Examples.
	//
	//    ;
	EmptyStmt struct {
		// Position of semicolon `;`.
		Semicolon int
	}

	// An ExprStmt node represents a stand-alone expression in a statement list.
	//
	// Examples.
	//
	//    42;
	//    f();
	ExprStmt struct {
		// Stand-alone expression.
		X Expr
	}

	// An IfStmt node represents an if statement.
	//
	// Examples.
	//
	//    if (x != 0) { x++; }
	//    if (i < max) { i; } else { max; }
	IfStmt struct {
		// Position of `if` keyword.
		If int
		// Condition.
		Cond Expr
		// True branch.
		Body Stmt
		// False branch; or nil if 1-way conditional.
		Else Stmt
	}

	// A ReturnStmt node represents a return statement.
	//
	// Examples.
	//
	//    return;
	//    return 42;
	ReturnStmt struct {
		// Position of `return` keyword.
		Return int
		// Result expression; or nil if void return.
		Result Expr
	}

	// A WhileStmt node represents a while statement.
	//
	// Examples.
	//
	//    while (i < 10) { i++; }
	WhileStmt struct {
		// Position of `while` keyword.
		While int
		// Condition.
		Cond Expr
		// Loop body.
		Body Stmt
	}
)

// A BlockItem represents an item of a block statement, and has one of the
// following underlying types.
//
//    Decl
//    Stmt
type BlockItem interface {
	Node
	// isBlockItem ensures that only block item nodes can be assigned to the
	// BlockItem interface.
	isBlockItem()
}

// An Expr node represents an expression, and has one of the following
// underlying types.
//
//    *BasicLit
//    *BinaryExpr
//    *CallExpr
//    *Ident
//    *IndexExpr
//    *ParenExpr
//    *UnaryExpr
type Expr interface {
	Node
	// isExpr ensures that only expression nodes can be assigned to the Expr
	// interface.
	isExpr()
}

// Expression nodes.
type (
	// A BasicLit node represents a basic literal.
	//
	// Examples.
	//
	//    42
	//    'a'
	BasicLit struct {
		// Position of basic literal.
		ValPos int
		// Basic literal type, one of the following.
		//
		//    token.CharLit
		//    token.IntLit
		Kind token.Kind
		// Basic literal value; e.g. 123, 'a'.
		Val string
	}

	// An BinaryExpr node represents a binary expression; X op Y.
	//
	// Examples.
	//
	//    x + y
	//    x = 42
	BinaryExpr struct {
		// First operand.
		X Expr
		// Position of operator.
		OpPos int
		// Operator, one of the following.
		//    token.Add      // +
		//    token.Sub      // -
		//    token.Mul      // *
		//    token.Div      // /
		//    token.Lt       // <
		//    token.Gt       // >
		//    token.Le       // <=
		//    token.Ge       // >=
		//    token.Ne       // !=
		//    token.Eq       // ==
		//    token.Land     // &&
		//    token.Assign   // =
		Op token.Kind
		// Second operand.
		Y Expr
	}

	// A CallExpr node represents a call expression.
	//
	// Examples.
	//
	//    foo()
	//    bar(42)
	CallExpr struct {
		// Function name.
		Name *Ident
		// Position of left-parenthesis `(`.
		Lparen int
		// Function arguments.
		Args []Expr
		// Position of right-parenthesis `)`.
		Rparen int
	}

	// An Ident node represents an identifier.
	//
	// Examples.
	//
	//    x
	//    int
	Ident struct {
		// Position of identifier.
		NamePos int
		// Identifier name.
		Name string
		// Corresponding function, variable or type definition. The declaration
		// mapping is added during the semantic analysis phase, based on the
		// lexical scope of the identifier.
		Decl Decl
	}

	// An IndexExpr node represents an array index expression.
	//
	// Examples.
	//
	//    buf[i]
	IndexExpr struct {
		// Array name.
		Name *Ident
		// Position of left-bracket `[`.
		Lbracket int
		// Array index.
		Index Expr
		// Position of right-bracket `]`.
		Rbracket int
	}

	// A ParenExpr node represents a parenthesised expression.
	ParenExpr struct {
		// Position of left-parenthesis `(`.
		Lparen int
		// Parenthesised expression.
		X Expr
		// Position of right-parenthesis `)`.
		Rparen int
	}

	// An UnaryExpr node represents an unary expression; op X.
	//
	// Examples.
	//
	//    -42
	//    !(x == 3 || x == 10)
	UnaryExpr struct {
		// Position of unary operator.
		OpPos int
		// Operator, one of the following.
		//    token.Sub   // -
		//    token.Not   // !
		Op token.Kind
		// Operand.
		X Expr
	}
)

// A Type node represents a type of µC, and has one of the following underlying
// types.
//
//    *ArrayType
//    *FuncType
//    *Ident
type Type interface {
	Node
	// isType ensures that only type nodes can be assigned to the Type interface.
	isType()
}

// Type nodes.
type (
	// An ArrayType node represents an array type.
	//
	// Examples.
	//
	//    int[]
	//    char[128]
	ArrayType struct {
		// Element type.
		Elem Type
		// Position of left-bracket `[`.
		Lbracket int
		// Array length.
		Len int
		// Position of right-bracket `]`.
		Rbracket int
	}

	// A FuncType node represents a function signature.
	//
	// Examples.
	//
	//    int(void)
	//    int(int a, int b)
	FuncType struct {
		// Return type.
		Result Type
		// Position of left-parenthesis `(`.
		Lparen int
		// Function parameters.
		Params []*VarDecl
		// Position of right-parenthesis `)`.
		Rparen int
	}
)

func (n *ArrayType) String() string {
	if n.Len > 0 {
		return fmt.Sprintf("%v[%d]", n.Elem, n.Len)
	}
	return fmt.Sprintf("%v[]", n.Elem)
}

func (n *BasicLit) String() string {
	return n.Val
}

func (n *BinaryExpr) String() string {
	// TODO: Verify that n.Op prints as "+" rather than "Add"
	return fmt.Sprintf("%v %v %v", n.X, n.Op, n.Y)
}

func (n *BlockStmt) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("{\n")
	for _, item := range n.Items {
		buf.WriteString(item.String() + "\n")
	}
	buf.WriteString("}\n")
	return buf.String()
}

func (n *CallExpr) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(n.Name.String())
	buf.WriteString("(")
	for i, arg := range n.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (n *EmptyStmt) String() string {
	return ";"
}

func (n *ExprStmt) String() string {
	return fmt.Sprintf("%v;", n.X)
}

func (n *File) String() string {
	buf := new(bytes.Buffer)
	for _, decl := range n.Decls {
		buf.WriteString(decl.String())
	}
	return buf.String()
}

func (n *FuncDecl) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(n.FuncType.Result.String())
	buf.WriteString(" ")
	buf.WriteString(n.FuncName.String())
	buf.WriteString("(")
	for i, param := range n.FuncType.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.VarType.String())
		buf.WriteString(" ")
		buf.WriteString(param.VarName.String())
	}
	buf.WriteString(")")
	// TODO: Decide if we want to print definition body. We do not want to do
	// it for most error messages. Use Print("%q %q", n, n.Body)? Same for
	// others?
	return buf.String()
}

func (n *FuncType) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(n.Result.String())
	buf.WriteString("(")
	for i, param := range n.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.VarType.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (n *Ident) String() string {
	return n.Name
}

func (n *IfStmt) String() string {
	if n.Else != nil {
		return fmt.Sprintf("if (%v) %v else %v", n.Cond, n.Body, n.Else)
	}
	return fmt.Sprintf("if (%v) %v", n.Cond, n.Body)
}

func (n *IndexExpr) String() string {
	return fmt.Sprintf("%v[%v]", n.Name, n.Index)
}

func (n *ParenExpr) String() string {
	return fmt.Sprintf("(%v)", n.X)
}

func (n *ReturnStmt) String() string {
	if n.Result != nil {
		return fmt.Sprintf("return %v;", n.Result)
	}
	return "return;"
}

func (n *TypeDef) String() string {
	return fmt.Sprintf("typedef %v %v;", n.DeclType, n.TypeName)
}

func (n *UnaryExpr) String() string {
	return fmt.Sprintf("%v%v", n.Op, n.X)
}

func (n *VarDecl) String() string {
	switch typ := n.VarType.(type) {
	case *ArrayType:
		if typ.Len > 0 {
			return fmt.Sprintf("%v %v[%d];", typ.Elem, n.VarName, typ.Len)
		}
		return fmt.Sprintf("%v %v[];", typ.Elem, n.VarName)
	default:
		return fmt.Sprintf("%v %v;", typ, n.VarName)
	}
}

func (n *WhileStmt) String() string {
	return fmt.Sprintf("while (%v) %v", n.Cond, n.Body)
}

// Start returns the start position of the node within the input stream.
func (n *ArrayType) Start() int {
	return n.Elem.Start()
}

// Start returns the start position of the node within the input stream.
func (n *BasicLit) Start() int {
	return n.ValPos
}

// Start returns the start position of the node within the input stream.
func (n *BinaryExpr) Start() int {
	return n.X.Start()
}

// Start returns the start position of the node within the input stream.
func (n *BlockStmt) Start() int {
	return n.Lbrace
}

// Start returns the start position of the node within the input stream.
func (n *CallExpr) Start() int {
	return n.Name.Start()
}

// Start returns the start position of the node within the input stream.
func (n *EmptyStmt) Start() int {
	return n.Semicolon
}

// Start returns the start position of the node within the input stream.
func (n *ExprStmt) Start() int {
	return n.X.Start()
}

// Start returns the start position of the node within the input stream.
func (n *File) Start() int {
	if len(n.Decls) > 0 {
		return n.Decls[0].Start()
	}
	return 0
}

// Start returns the start position of the node within the input stream.
func (n *FuncDecl) Start() int {
	return n.FuncType.Start()
}

// Start returns the start position of the node within the input stream.
func (n *FuncType) Start() int {
	return n.Result.Start()
}

// Start returns the start position of the node within the input stream.
func (n *Ident) Start() int {
	return n.NamePos
}

// Start returns the start position of the node within the input stream.
func (n *IfStmt) Start() int {
	return n.If
}

// Start returns the start position of the node within the input stream.
func (n *IndexExpr) Start() int {
	return n.Name.Start()
}

// Start returns the start position of the node within the input stream.
func (n *ParenExpr) Start() int {
	return n.Lparen
}

// Start returns the start position of the node within the input stream.
func (n *ReturnStmt) Start() int {
	return n.Return
}

// Start returns the start position of the node within the input stream.
func (n *TypeDef) Start() int {
	return n.Typedef
}

// Start returns the start position of the node within the input stream.
func (n *UnaryExpr) Start() int {
	return n.OpPos
}

// Start returns the start position of the node within the input stream.
func (n *VarDecl) Start() int {
	return n.VarType.Start()
}

// Start returns the start position of the node within the input stream.
func (n *WhileStmt) Start() int {
	return n.While
}

// Verify that all nodes implement the Node interface.
var (
	_ Node = &ArrayType{}
	_ Node = &BasicLit{}
	_ Node = &BinaryExpr{}
	_ Node = &BlockStmt{}
	_ Node = &CallExpr{}
	_ Node = &EmptyStmt{}
	_ Node = &ExprStmt{}
	_ Node = &File{}
	_ Node = &FuncDecl{}
	_ Node = &FuncType{}
	_ Node = &Ident{}
	_ Node = &IfStmt{}
	_ Node = &IndexExpr{}
	_ Node = &ParenExpr{}
	_ Node = &ReturnStmt{}
	_ Node = &TypeDef{}
	_ Node = &UnaryExpr{}
	_ Node = &VarDecl{}
	_ Node = &WhileStmt{}
)

// Type returns the type of the declared identifier.
func (n *FuncDecl) Type() types.Type {
	// TODO: Consider caching the types.Type.
	return newType(n.FuncType)
}

// Type returns the type of the declared identifier.
func (n *VarDecl) Type() types.Type {
	// TODO: Consider caching the types.Type.
	return newType(n.VarType)
}

// Type returns the type of the declared identifier.
func (n *TypeDef) Type() types.Type {
	// TODO: Consider caching the types.Type.

	// NOTE: "A typedef declaration does not introduce a new type, only a synonym
	// for the type so specified." (see §6.7.7.3)
	return newType(n.DeclType)
}

// Name returns the name of the declared identifier.
func (n *FuncDecl) Name() *Ident {
	return n.FuncName
}

// Name returns the name of the declared identifier.
func (n *VarDecl) Name() *Ident {
	return n.VarName
}

// Name returns the name of the declared identifier.
func (n *TypeDef) Name() *Ident {
	return n.TypeName
}

// Value returns the initializing value of the defined identifier; or nil if
// declaration or tentative definition.
//
// Underlying type for function declarations.
//
//    *BlockStmt
func (n *FuncDecl) Value() Node {
	return n.Body
}

// Value returns the initializing value of the defined identifier; or nil if
// declaration or tentative definition.
//
// Underlying type for variable declarations.
//
//    Expr
func (n *VarDecl) Value() Node {
	return n.Val
}

// Value returns the initializing value of the defined identifier; or nil if
// declaration or tentative definition.
//
// Underlying type for type definitions.
//
//    Type
func (n *TypeDef) Value() Node {
	return n.DeclType
}

// isDecl ensures that only declaration nodes can be assigned to the Decl
// interface.
func (n *FuncDecl) isDecl() {}
func (n *VarDecl) isDecl()  {}
func (n *TypeDef) isDecl()  {}

// Verify that the declaration nodes implement the Decl interface.
var (
	_ Decl = &FuncDecl{}
	_ Decl = &VarDecl{}
	_ Decl = &TypeDef{}
)

// isStmt ensures that only statement nodes can be assigned to the Stmt
// interface.
func (n *BlockStmt) isStmt()  {}
func (n *EmptyStmt) isStmt()  {}
func (n *ExprStmt) isStmt()   {}
func (n *IfStmt) isStmt()     {}
func (n *ReturnStmt) isStmt() {}
func (n *WhileStmt) isStmt()  {}

// Verify that the statement nodes implement the Stmt interface.
var (
	_ Stmt = &BlockStmt{}
	_ Stmt = &EmptyStmt{}
	_ Stmt = &ExprStmt{}
	_ Stmt = &IfStmt{}
	_ Stmt = &ReturnStmt{}
	_ Stmt = &WhileStmt{}
)

// isBlockItem ensures that only block item nodes can be assigned to the
// BlockItem interface.
func (n *BlockStmt) isBlockItem()  {}
func (n *EmptyStmt) isBlockItem()  {}
func (n *ExprStmt) isBlockItem()   {}
func (n *FuncDecl) isBlockItem()   {}
func (n *IfStmt) isBlockItem()     {}
func (n *ReturnStmt) isBlockItem() {}
func (n *TypeDef) isBlockItem()    {}
func (n *VarDecl) isBlockItem()    {}
func (n *WhileStmt) isBlockItem()  {}

// Verify that the block item nodes implement the BlockItem interface.
var (
	_ BlockItem = &BlockStmt{}
	_ BlockItem = &EmptyStmt{}
	_ BlockItem = &ExprStmt{}
	_ BlockItem = &FuncDecl{}
	_ BlockItem = &IfStmt{}
	_ BlockItem = &ReturnStmt{}
	_ BlockItem = &TypeDef{}
	_ BlockItem = &VarDecl{}
	_ BlockItem = &WhileStmt{}
)

// isExpr ensures that only expression nodes can be assigned to the Expr
// interface.
func (n *BasicLit) isExpr()   {}
func (n *BinaryExpr) isExpr() {}
func (n *CallExpr) isExpr()   {}
func (n *Ident) isExpr()      {}
func (n *IndexExpr) isExpr()  {}
func (n *ParenExpr) isExpr()  {}
func (n *UnaryExpr) isExpr()  {}

// Verify that the expression nodes implement the Expr interface.
var (
	_ Expr = &BasicLit{}
	_ Expr = &BinaryExpr{}
	_ Expr = &CallExpr{}
	_ Expr = &Ident{}
	_ Expr = &IndexExpr{}
	_ Expr = &ParenExpr{}
	_ Expr = &UnaryExpr{}
)

// isType ensures that only type nodes can be assigned to the Type interface.
func (n *Ident) isType()     {}
func (n *ArrayType) isType() {}
func (n *FuncType) isType()  {}

// Verify that the type nodes implement the Type interface.
var (
	_ Type = &Ident{}
	_ Type = &ArrayType{}
	_ Type = &FuncType{}
)
