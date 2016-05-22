package irgen

import (
	"fmt"
	"log"
	"strconv"

	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/instruction"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/term"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/token"
	uctypes "github.com/mewmew/uc/types"
)

// generator hold the variables needed for recursing over node types.
type generator struct {
	// module is the current module being created
	module ir.Module
	// funcDefStack holds all current function nests
	funcDefStack []*ir.Function
	// currentFunction points to the function being created
	currentFunction *ir.Function
	// basicBlocks holds basic blocks before function creation
	basicBlocks []*ir.BasicBlock
	// instructionBuffer is appens instruction before basic block creation.
	instructionBuffer []instruction.Instruction
	// ssaCounter give a unique id for anonymous assignments and basic blocks.
	ssaCounter int
	// lastLabel holds the ssa of the last created basic block
	lastLabel int
}

// recurse calls the appropriate function for the type of node.
func (gen *generator) recurse(n ast.Node) error {
	var err error
	switch n := n.(type) {
	case *ast.File:
		err = gen.createFile(n)
	case *ast.BlockStmt:
		err = gen.createBlock(n)
	case *ast.FuncDecl:
		err = gen.createFunction(n)
	case *ast.CallExpr:
		err = gen.createCall(n)
	case *ast.VarDecl:
		err = gen.createVar(n)
	case *ast.Ident:
		err = gen.loadIdent(n)
	case *ast.BinaryExpr:
		err = gen.createBinaryExpr(n)
	case *ast.ExprStmt:
		err = gen.createExprStmt(n)
	case *ast.WhileStmt:
		err = gen.createWhile(n)
	default:
		log.Printf("Not yet implemented: type %#v", n)
	}
	// TODO: Implement the rest of the needed node types
	if err != nil {
		return errutil.Err(err)
	}
	return nil
}

// newGenerator creates a new generator object.
func newGenerator() *generator {
	gen := new(generator)
	// The usual size of a basic block seems to be less than 10 instructions.
	gen.instructionBuffer = make([]instruction.Instruction, 0, 10)
	gen.ssaCounter = 0
	gen.lastLabel = gen.ssaCounter
	return gen
}

// Gen generates LLVM IR based on the syntax tree of the given file.
func Gen(file *ast.File) error {
	// TODO: REMOVE log messages
	log.SetPrefix(term.BlueBold("Log:"))
	log.SetFlags(log.Lshortfile)
	gen := newGenerator()
	if err := gen.recurse(file); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// createFile recurses over the declarations in an *ast.File.
func (gen *generator) createFile(n *ast.File) error {
	for _, decl := range n.Decls {
		if err := gen.recurse(decl); err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// createBlock creates the instructions for a block statement.
func (gen *generator) createBlock(n *ast.BlockStmt) error {
	for _, blockItem := range n.Items {
		log.Printf("blockItem %#v of type %T\n", blockItem, blockItem)
		if err := gen.recurse(blockItem); err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// createFunction creates the instructions for creating a function.
func (gen *generator) createFunction(n *ast.FuncDecl) error {
	//var fn *ir.Function
	//fn, ssaCounter = createFunction(n, ssaCounter)

	// TODO: Create unique names for nested functions, this should
	// really be done in parser... the following currently breakes
	// ability to call nested functions
	var name string
	for _, prevfn := range gen.funcDefStack {
		name += prevfn.Name() + "."
	}
	name += n.Name().String()
	sig := toIrType(n.Type())
	var fn *ir.Function
	if sig, ok := sig.(*irtypes.Func); ok {
		fn = ir.NewFunction(name, sig)
	} else {
		panic(errutil.Newf("Conversion from uc function type failed: %#v %#v", sig, n.Type()))
	}
	gen.module.Funcs = append(gen.module.Funcs, fn)

	if !astutil.IsDef(n) {
		log.Printf("create function decl %v\n", fn)
		return nil
	}

	log.Printf("create function def %v\n", fn)
	gen.funcDefStack = append(gen.funcDefStack, fn)
	gen.currentFunction = fn

	// Recurse body
	if err := gen.recurse(n.Body); err != nil {
		return errutil.Err(err)
	}
	ret, err := instruction.NewRet(irtypes.NewVoid(), nil)
	if err != nil {
		return errutil.Err(err)
	}
	allInsts := make([]instruction.Instruction, len(gen.instructionBuffer))
	copy(allInsts, gen.instructionBuffer)
	gen.instructionBuffer = gen.instructionBuffer[:0]
	bb, err := ir.NewBasicBlock(encLocal(gen.lastLabel), allInsts, ret)
	if err != nil {
		return errutil.Err(err)
	}
	gen.basicBlocks = append(gen.basicBlocks, bb)
	log.Print(gen.basicBlocks)
	gen.ssaCounter++
	gen.funcDefStack = gen.funcDefStack[:len(gen.funcDefStack)-1]
	if len(gen.funcDefStack) > 0 {
		gen.currentFunction = gen.funcDefStack[len(gen.funcDefStack)-1]
	} else {
		gen.currentFunction = nil
	}
	return nil
}

// createCall creates the instructions for calling a function.
func (gen *generator) createCall(n *ast.CallExpr) error {
	var insts []instruction.Instruction
	log.Printf("%v: create call %v with:\n", encLocal(gen.ssaCounter), n)
	if callType, ok := n.Name.Decl.Type().(*uctypes.Func); ok {
		// TODO: Add support for elipsis/variadic functions
		for i, arg := range n.Args {
			log.Printf("Arg %v of type %v", arg.String(), callType.Params[i].Type.String())
		}
		//instruction.NewCall(toIrType(callType.Result), n.Name.String(), args)
		gen.ssaCounter++
		//instruction.NewLocalVarDef("")
	} else {
		return errutil.Newf("invalid type assertion; expected: *types.Func, got: %T", n.Name.Decl.Type())
	}
	gen.instructionBuffer = append(gen.instructionBuffer, insts...)
	return nil
}

// createVar creates the instruction for allocating place for a variable.
func (gen *generator) createVar(n *ast.VarDecl) error {
	// Protect against void parameter
	// TODO: Check that only those are caught
	if uctypes.IsVoid(n.Type()) {
		log.Printf("Void caught, hopfully in function definition parameter: %#v", n)
		return nil
	}
	if len(gen.funcDefStack) > 0 {
		err := gen.createLocal(n)
		if err != nil {
			return errutil.Err(err)
		}
	} else {
		// Global values are compile time constant, no need for ssa
		// NOTE: Global variables may still be unnamed, but they would have
		// a different counter; e.g. @0 and %0 may co-exist. We may split
		// ssaCounter into a local and a global counter, along the lines of:
		//
		//    localCounter := 0
		//    globalCounter := 0
		err := gen.createGlobal(n)
		if err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// loadIdent creates the instruction for loading an ident into a local var.
func (gen *generator) loadIdent(n *ast.Ident) error {
	gen.ssaCounter++
	log.Printf("loading ident %v into %v", n.Name, encLocal(gen.ssaCounter))
	typ := toIrType(n.Decl.Type())
	var newVar value.Value
	var err error
	// TODO: add data fields or deduce if n is local or global

	ptrType, err := irtypes.NewPointer(typ)
	if err != nil {
		return errutil.Err(err)
	}
	if true {
		newVar, err = ir.NewLocal(n.Name, ptrType)
	} else {
		newVar, err = constant.NewPointer(ptrType, n.Name)
	}
	if err != nil {
		return errutil.Err(err)
	}
	lvd, err := instruction.NewLoad(typ, newVar)
	if err != nil {
		return errutil.Err(err)
	}
	instr, err := instruction.NewLocalVarDef("", lvd)
	if err != nil {
		return errutil.Err(err)
	}
	gen.instructionBuffer = append(gen.instructionBuffer, instr)
	return nil
}

// createBinaryExpr creates the instructions for a binary expression.
func (gen *generator) createBinaryExpr(n *ast.BinaryExpr) error {
	log.Println("bin expr")
	log.Println(n)
	if err := gen.recurse(n.Y); err != nil {
		return errutil.Err(err)
	}
	_ = gen.ssaCounter
	switch n.Op {
	case token.Assign:
		log.Println(n.X)
		if x, ok := n.X.(*ast.Ident); ok {
			log.Printf("Assign %v %v %v to %#v", n.X, n.Op, n.Y, x)
			return nil
		}
		//default:
	}
	gen.ssaCounter++
	log.Printf("Assign %v %v %v to %#v", n.X, n.Op, n.Y, encLocal(gen.ssaCounter))
	return nil
}

// createExprStmt creates the instructions for an expression statment.
func (gen *generator) createExprStmt(n *ast.ExprStmt) error {
	if err := gen.recurse(n.X); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// createLocal creates the instructions for a local variable declaration.
func (gen *generator) createLocal(n *ast.VarDecl) error {
	// TODO: Implement
	var insts []instruction.Instruction
	log.Printf("create local variable %v\n", n)

	gen.instructionBuffer = append(gen.instructionBuffer, insts...)
	return nil
}

// createGlobal creates the instructions for a global variable declaration.
func (gen *generator) createGlobal(n *ast.VarDecl) error {
	// TODO: Implement
	log.Printf("create global variable %v\n", n)
	var gv *ir.GlobalDecl
	gen.module.Globals = append(gen.module.Globals, gv)
	return nil
}

// createWhile creates the instructions for a while loop.
func (gen *generator) createWhile(n *ast.WhileStmt) error {
	log.Printf("start while loop %v\n", n)
	// Finnish last basic block with a branch to the basic block we
	// are about to create (with label ssaCounter+1)
	gen.ssaCounter++
	whileLabel := gen.ssaCounter
	allInsts := make([]instruction.Instruction, len(gen.instructionBuffer))
	copy(allInsts, gen.instructionBuffer)
	log.Printf("Clear last basic block instrucitons: %v\n", allInsts)
	brToWhileLabel, err := instruction.NewBr(nil, encLocal(whileLabel), "")
	if err != nil {
		return errutil.Err(err)
	}
	//terminator, gen.ssaCounter = createWhile(n, gen.ssaCounter)
	bb, err := ir.NewBasicBlock(encLocal(gen.lastLabel), allInsts, brToWhileLabel)
	if err != nil {
		return errutil.Err(err)
	}
	gen.basicBlocks = append(gen.basicBlocks, bb)

	gen.ssaCounter++
	gen.lastLabel = gen.ssaCounter
	gen.instructionBuffer = gen.instructionBuffer[:0]

	// Recurse over body of while loop
	log.Printf("while.Body = %#v\n", n.Body)

	if err := gen.recurse(n.Body); err != nil {
		return errutil.Err(err)
	}

	//_, gen.ssaCounter = endWhile(n, gen.ssaCounter)

	// End the while loop with a branch to whileLabel
	bb, err = ir.NewBasicBlock(encLocal(gen.lastLabel), allInsts, brToWhileLabel)
	if err != nil {
		return errutil.Err(err)
	}
	gen.basicBlocks = append(gen.basicBlocks, bb)
	gen.ssaCounter++
	gen.lastLabel = gen.ssaCounter
	gen.instructionBuffer = gen.instructionBuffer[:0]
	// TODO: Implement

	return nil
}

// encLocal makes a anonymous local variable string representation from int.
func encLocal(ssa int) string {
	return asm.EncLocal(strconv.Itoa(ssa))
}

func toIrType(n uctypes.Type) irtypes.Type {
	//TODO: implement, placeholder implementation
	var t irtypes.Type
	var err error
	switch ucType := n.(type) {
	case *uctypes.Basic:
		switch ucType.Kind {
		case uctypes.Int:
			//TODO: Get int width from compile env
			t, err = irtypes.NewInt(32)
		case uctypes.Char:
			t, err = irtypes.NewInt(8)
		case uctypes.Void:
			t = irtypes.NewVoid()
		}
	case *uctypes.Array:
		elem := toIrType(ucType.Elem)
		t, err = irtypes.NewArray(elem, ucType.Len)
	case *uctypes.Func:
		var params []*irtypes.Param
		variadic := false
		for _, p := range ucType.Params {
			//TODO: Add support for variadic
			if uctypes.IsVoid(p.Type) {
				break
			}
			pt := toIrType(p.Type)
			log.Printf("converting type %#v to %#v", p.Type, pt)
			params = append(params, irtypes.NewParam(pt, ""))
		}
		result := toIrType(ucType.Result)
		t, err = irtypes.NewFunc(result, params, variadic)
	default:
		panic(fmt.Sprintf("support for translating type %T not yet implemented.", ucType))
	}
	if err != nil {
		panic(errutil.Err(err))
	}
	if t == nil {
		panic(errutil.Newf("Conversion failed: %#v\n", n))
	}
	return t
}
