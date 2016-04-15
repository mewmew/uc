// TODO: Add source position tracking of nodes.

// Package ast declares the types used to represent abstract syntax trees of µC
// soure code.
package ast

import (
	"github.com/mewmew/uc/token"
	"github.com/mewmew/uc/types"
)

// A File represents a µC source file.
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
type Node interface {
	// Start returns the start position of the node within the input stream.
	Start() int
}

// A Decl node represents a declaration, and has one of the following underlying
// types.
//
//    *FuncDecl
//    *VarDecl
type Decl interface {
	Node
	isDecl()
}

// Declaration nodes.
type (
	// A FuncDecl node represents a function declaration.
	//
	// Examples.
	//
	//    int add(int a, int b) { return a+b; }
	//    int puts(char s[]);
	FuncDecl struct {
		// Function signature.
		Type *types.Func
		// Function name.
		Name *Ident
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
		Type types.Type
		// Variable name.
		Name *Ident
		// Variable value expression; or nil if variable declaration (i.e. not
		// variable definition).
		Val Expr
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
		// List of block items contained within the block.
		Items []BlockItem
	}

	// An EmptyStmt node represents an empty statement (i.e. ";").
	EmptyStmt struct{}

	// An ExprStmt node represents a stand-alone expression in a statement list.
	ExprStmt struct {
		// Stand-alone expression.
		X Expr
	}

	// An IfStmt node represents an if statement.
	IfStmt struct {
		// Condition.
		Cond Expr
		// True branch.
		Body Stmt
		// False branch; or nil if 1-way conditional.
		Else Stmt
	}

	// A WhileStmt node represents a while statement.
	WhileStmt struct {
		// Condition.
		Cond Expr
		// Loop body.
		Body Stmt
	}

	// A ReturnStmt node represents a return statement.
	ReturnStmt struct {
		// Result expression; or nil if void return.
		Result Expr
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
	// An Ident node represents an identifier.
	Ident struct {
		// Identifier name.
		Name string
	}

	// A BasicLit node represents a basic literal.
	BasicLit struct {
		// Basic literal type, one of the following.
		//
		//    token.CharLit
		//    token.IntLit
		Kind token.Kind
		// Basic literal value; e.g. 123, 'a'.
		Val string
	}

	// An UnaryExpr node represents an unary expression; op X.
	UnaryExpr struct {
		// Operator, one of the following.
		//    token.Sub   // -
		//    token.Not   // !
		Op token.Kind
		// Operand.
		X Expr
	}

	// An BinaryExpr node represents a binary expression; X op Y.
	BinaryExpr struct {
		// First operand.
		X Expr
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
	CallExpr struct {
		// Function name.
		Name *Ident
		// Function arguments.
		Args []Expr
	}

	// A ParenExpr node represents a parenthesised expression.
	ParenExpr struct {
		// Parenthesised expression.
		X Expr
	}

	// An IndexExpr node represents an array index expression.
	IndexExpr struct {
		// Array name.
		Name *Ident
		// Array index.
		Index Expr
	}
)

// Start returns the start position of the node within the input stream.
func (n *BasicLit) Start() int { panic("ast.BasicLit.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *BinaryExpr) Start() int { panic("ast.BinaryExpr.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *BlockStmt) Start() int { panic("ast.BlockStmt.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *CallExpr) Start() int { panic("ast.CallExpr.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *EmptyStmt) Start() int { panic("ast.EmptyStmt.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *ExprStmt) Start() int { panic("ast.ExprStmt.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *File) Start() int { panic("ast.File.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *FuncDecl) Start() int { panic("ast.FuncDecl.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *Ident) Start() int { panic("ast.Ident.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *IfStmt) Start() int { panic("ast.IfStmt.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *IndexExpr) Start() int { panic("ast.IndexExpr.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *ParenExpr) Start() int { panic("ast.ParenExpr.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *ReturnStmt) Start() int { panic("ast.ReturnStmt.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *UnaryExpr) Start() int { panic("ast.UnaryExpr.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *VarDecl) Start() int { panic("ast.VarDecl.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *WhileStmt) Start() int { panic("ast.WhileStmt.Start: not yet implemented") }

// Verify that all nodes implement the Node interface.
var (
	_ Node = &BasicLit{}
	_ Node = &BinaryExpr{}
	_ Node = &BlockStmt{}
	_ Node = &CallExpr{}
	_ Node = &EmptyStmt{}
	_ Node = &ExprStmt{}
	_ Node = &File{}
	_ Node = &FuncDecl{}
	_ Node = &Ident{}
	_ Node = &IfStmt{}
	_ Node = &IndexExpr{}
	_ Node = &ParenExpr{}
	_ Node = &ReturnStmt{}
	_ Node = &UnaryExpr{}
	_ Node = &VarDecl{}
	_ Node = &WhileStmt{}
)

// isDecl ensures that only declaration nodes can be assigned to the Decl
// interface.
func (n *FuncDecl) isDecl() {}
func (n *VarDecl) isDecl()  {}

// Verify that the declaration nodes implement the Decl interface.
var (
	_ Decl = &FuncDecl{}
	_ Decl = &VarDecl{}
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
