// Package ast declares the types used to represent abstract syntax trees of µC
// soure code.
package ast

import "github.com/mewmew/uc/token"

// A File represents a µC source file.
type File struct {
	// Top-level declarations.
	Decls []TopLevelDecl
}

// A Node represents a node within the abstract syntax tree, and has one of the
// following underlying types.
//
//    TopLevelDecl
//    Decl
//    Stmt
//    Expr
type Node interface {
	// Start returns the start position of the node within the input stream.
	Start() int
	// End returns the first character immediately after the node within the
	// input stream.
	End() int
}

// TODO: Evaluate whether TopLevelDecl should be merged with Decl, to simplify
// the structure of the AST; in which case a semantic analysis pass would be
// added to ensure that function declarations are not nested (even if this would
// be an interesting extension to C).

// A TopLevelDecl node represents a top-level declaration, and has one of the
// following underlying types.
//
//    Decl
//    *FuncDecl
type TopLevelDecl interface {
	Node
	// isTopLevelDecl ensures that only top-level declaration nodes can be
	// assigned to the TopLevelDevl interface.
	isTopLevelDecl()
}

// Top-level declaration nodes.
type (
	// A FuncDecl node represents a function declaration.
	FuncDecl struct {
		// Function name.
		Name *Ident
		// Function signature.
		Type *FuncType
		// Function body; or nil if function declaration (i.e. not function
		// definition).
		Body *BlockStmt
	}
)

// A Decl node represents a declaration, and has one of the following underlying
// types.
//
//    *VarDecl
type Decl interface {
	Node
	// isDecl ensures that only declaration nodes can be assigned to the Decl
	// interface.
	isDecl()
}

// Declaration nodes.
type (
	// A VarDecl node represents a variable declaration.
	VarDecl struct {
		// Variable name.
		Name *Ident
		// Variable type.
		Type *VarType
		// Variable value expression; or nil if variable declaration (i.e. not
		// variable definition).
		Val Expr
	}
)

// A Stmt node represents a statement, and has one of the following underlying
// types.
//
//    *BlockStmt
//    TODO
type Stmt interface {
	Node
	// isStmt ensures that only statement nodes can be assigned to the Stmt
	// interface.
	isStmt()
}

// TODO: Add semantic analysis pass which verifies that declaration statements
// precedes any other statements in the body of function block.
//
//    // first specifies the first non-declaration statement within the
//    // statements of the block.
//    first := -1
//    for i, stmt := f.Body.Stmts {
//       if _, ok := stmt.(*DeclStmt); ok {
//          if first != -1 {
//             return errutil.Newf("declaration statement %v occurs after first non-declaration statement %v in function body", stmt, f.Body.Stmts[first])
//          }
//       } else if first == -1 {
//          first = i
//       }
//    }

// Statement nodes.
type (
	// A BlockStmt node represents a block statement.
	BlockStmt struct {
		// List of statements contained within the block.
		Stmts []Stmt
	}

	// A DeclStmt node represents a declaration statement.
	DeclStmt struct {
		// TODO: Add support for lists of declarations? E.g.
		//    int x, y, z;
		//    int x, y = 3, *z;

		// Declaration.
		Decl Decl
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

// An Expr node represents an expression, and has one of the following
// underlying types.
//
//    TODO
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

	// An UnaryExpr node represents an unary expression; op B.
	UnaryExpr struct {
		// Operator, one of the following.
		//    token.Sub   // -
		//    token.Not   // !
		Op token.Kind
		// Operand.
		B Expr
	}

	// An BinaryExpr node represents a binary expression; A op B.
	BinaryExpr struct {
		// First operand.
		A Expr
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
		B Expr
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

// TODO: Move FuncType and VarType to dedicated types package.
// TODO: Add definition of FuncType and VarType.

// A FuncType represents a type signature.
type FuncType struct{}

// A VarType represents a variable type.
type VarType struct{}

// Start returns the start position of the node within the input stream.
func (n *BasicLit) Start() int { panic("ast.BasicLit.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *BinaryExpr) Start() int { panic("ast.BinaryExpr.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *BlockStmt) Start() int { panic("ast.BlockStmt.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *CallExpr) Start() int { panic("ast.CallExpr.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *DeclStmt) Start() int { panic("ast.DeclStmt.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *EmptyStmt) Start() int { panic("ast.EmptyStmt.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *ExprStmt) Start() int { panic("ast.ExprStmt.Start: not yet implemented") }

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

// End returns the first character immediately after the node within the input
// stream.
func (n *BasicLit) End() int { panic("ast.BasicLit.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *BinaryExpr) End() int { panic("ast.BinaryExpr.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *BlockStmt) End() int { panic("ast.BlockStmt.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *CallExpr) End() int { panic("ast.CallExpr.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *DeclStmt) End() int { panic("ast.DeclStmt.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *EmptyStmt) End() int { panic("ast.EmptyStmt.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *ExprStmt) End() int { panic("ast.ExprStmt.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *FuncDecl) End() int { panic("ast.FuncDecl.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *Ident) End() int { panic("ast.Ident.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *IfStmt) End() int { panic("ast.IfStmt.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *IndexExpr) End() int { panic("ast.IndexExpr.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *ParenExpr) End() int { panic("ast.ParenExpr.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *ReturnStmt) End() int { panic("ast.ReturnStmt.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *UnaryExpr) End() int { panic("ast.UnaryExpr.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *VarDecl) End() int { panic("ast.VarDecl.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *WhileStmt) End() int { panic("ast.WhileStmt.End: not yet implemented") }

// Verify that all nodes implement the Node interface.
var (
	_ Node = &BasicLit{}
	_ Node = &BinaryExpr{}
	_ Node = &BlockStmt{}
	_ Node = &CallExpr{}
	_ Node = &DeclStmt{}
	_ Node = &EmptyStmt{}
	_ Node = &ExprStmt{}
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

// isTopLevelDecl ensures that only top-level declaration nodes can be assigned
// to the TopLevelDevl interface.
func (n *FuncDecl) isTopLevelDecl() {}
func (n *VarDecl) isTopLevelDecl()  {}

// Verify that the top-level declaration nodes implement the TopLevelDecl
// interface.
var (
	_ TopLevelDecl = &FuncDecl{}
	_ TopLevelDecl = &VarDecl{}
)

// isDecl ensures that only declaration nodes can be assigned to the Decl
// interface.
func (n *VarDecl) isDecl() {}

// Verify that the declaration nodes implement the Decl interface.
var (
	_ Decl = &VarDecl{}
)

// isStmt ensures that only statement nodes can be assigned to the Stmt
// interface.
func (n *BlockStmt) isStmt()  {}
func (n *DeclStmt) isStmt()   {}
func (n *EmptyStmt) isStmt()  {}
func (n *ExprStmt) isStmt()   {}
func (n *IfStmt) isStmt()     {}
func (n *ReturnStmt) isStmt() {}
func (n *WhileStmt) isStmt()  {}

// Verify that the statement nodes implement the Stmt interface.
var (
	_ Stmt = &BlockStmt{}
	_ Stmt = &DeclStmt{}
	_ Stmt = &EmptyStmt{}
	_ Stmt = &ExprStmt{}
	_ Stmt = &IfStmt{}
	_ Stmt = &ReturnStmt{}
	_ Stmt = &WhileStmt{}
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
