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
	"github.com/mewmew/uc/sem"
	"github.com/mewmew/uc/token"
	uctypes "github.com/mewmew/uc/types"
)

// generator hold the variables needed for recursing over node types.
type generator struct {
	// module is the current module being created
	module *ir.Module
	// info holds usefull info from the semantic checker
	info *sem.Info
	// funcDefStack holds all current function nests
	funcDefStack []*ir.Function
	// currentFunction points to the function being created
	currentFunction *ir.Function
	// basicBlocks holds basic blocks before function creation
	basicBlocks []*ir.BasicBlock
	// instructionBuffer holds instructions before basic block creation.
	instructionBuffer []instruction.Instruction
	// ssaCounter give a unique id for anonymous assignments and basic blocks,
	// starting at 0.
	ssaCounter int
	// lastLabel holds the ssa of the last created basic block
	lastLabel int
}

// recurse calls the appropriate function for the type of node.
func (gen *generator) recurse(n ast.Node) error {
	var err error
	switch n := n.(type) {
	case *ast.BasicLit:
		err = gen.loadBasicLit(n)
	case *ast.BinaryExpr:
		err = gen.createBinaryExpr(n)
	case *ast.BlockStmt:
		err = gen.createBlock(n)
	case *ast.CallExpr:
		err = gen.createCall(n)
	case *ast.ExprStmt:
		err = gen.createExprStmt(n)
	case *ast.File:
		err = gen.createFile(n)
	case *ast.FuncDecl:
		err = gen.createFunction(n)
	case *ast.Ident:
		err = gen.loadIdent(n)
	case *ast.UnaryExpr:
		err = gen.createUnaryExpr(n)
	case *ast.VarDecl:
		err = gen.createVar(n)
	case *ast.WhileStmt:
		err = gen.createWhile(n)
	case *ast.EmptyStmt:
		// nothing to do.
		return nil
	case *ast.IndexExpr:
		err = gen.loadIndexExpr(n)
	default:
		log.Printf("recurse does not yet implement type %#v", n)
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
	gen.module = new(ir.Module)
	// The usual size of a basic block seems to be less than 10 instructions.
	gen.instructionBuffer = make([]instruction.Instruction, 0, 10)
	gen.lastLabel = gen.ssaCounter
	return gen
}

// Gen generates LLVM IR based on the syntax tree of the given file.
func Gen(file *ast.File, info *sem.Info) (*ir.Module, error) {
	// TODO: REMOVE log messages
	log.SetPrefix(term.BlueBold("Log:"))
	log.SetFlags(log.Lshortfile)
	gen := newGenerator()
	gen.info = info
	if err := gen.recurse(file); err != nil {
		return nil, errutil.Err(err)
	}
	return gen.module, nil
}

// loadBasicLit loads the basic lit into next ssaCount.
func (gen *generator) loadBasicLit(n *ast.BasicLit) error {
	gen.ssaCounter++
	//TODO: Fix early conversions, get type, now highest accuracy.
	typ, err := irtypes.NewInt(32)
	if err != nil {
		return errutil.Err(err)
	}
	con, err := constant.NewInt(typ, n.Val)
	if err != nil {
		return errutil.Err(err)
	}
	zero, err := constant.NewInt(typ, "0")
	if err != nil {
		return errutil.Err(err)
	}
	val, err := instruction.NewAdd(con, zero)
	if err != nil {
		return errutil.Err(err)
	}
	instr, err := instruction.NewLocalVarDef("", val)
	// Add 'add n.Val 0' instruction for storing in anon local var
	gen.instructionBuffer = append(gen.instructionBuffer, instr)
	return nil
}

// createBinaryExpr creates the instructions for a binary expression.
func (gen *generator) createBinaryExpr(n *ast.BinaryExpr) error {
	xtype := gen.info.Types[n.X]
	ytype := gen.info.Types[n.Y]
	log.Printf("bin expr: n.X: %v %v; Op: %v; n.Y: %v %v", toIrType(xtype), n.X, n.Op, toIrType(ytype), n.Y)
	if err := gen.recurse(n.Y); err != nil {
		return errutil.Err(err)
	}
	yssa := gen.ssaCounter
	yval, ok := gen.instructionBuffer[len(gen.instructionBuffer)-1].(value.Value)
	if !ok {
		panic("Last instruction not a value.Value")
	}

	// handle assignment specially from ident
	if n.Op == token.Assign {
		log.Printf("saving %%%v to %v", yssa, n.X)
		err := gen.store(n.X)
		if err != nil {
			return errutil.Err(err)
		}
		return nil
	}

	//not an assignment, evaluate both sides

	if err := gen.recurse(n.X); err != nil {
		return errutil.Err(err)
	}
	xssa := gen.ssaCounter
	xval, ok := gen.instructionBuffer[len(gen.instructionBuffer)-1].(value.Value)
	if !ok {
		panic("Last instruction not a value.Value")
	}

	log.Printf("Y: %v evaluated in %%%v", n.Y, yssa)
	log.Printf("Y: %v evaluated in %%%v", n.X, xssa)
	log.Printf("comparing %%%v to %%%v", xssa, yssa)

	switch n.Op {
	case token.Eq:
		log.Printf("with cond ICondEq")
		gen.ssaCounter++
		instruction.NewICmp(instruction.ICondEq, xval, yval)
		err := gen.store(n.X)
		if err != nil {
			return errutil.Err(err)
		}
	default:
		gen.ssaCounter++
		log.Printf("%v %v %v saved to %#v ", n.X, n.Op, n.Y, encLocal(gen.ssaCounter))
	}
	return nil
}

func (gen *generator) getElemPtr(addr value.Value, index ast.Expr) error {
	// The often misunderstood GEP instruction's most likely grotesquely
	// misguided implementation
	i64, err := irtypes.NewInt(64)
	if err != nil {
		return errutil.Err(err)
	}
	zero, err := constant.NewInt(i64, "0")
	indices := []value.Value{zero}
	switch n := index.(type) {
	case *ast.BasicLit:

		strconv.Atoi(n.Val)
		val, err := instruction.NewGetElementPtr(i64, addr, indices)
		if err != nil {
			return errutil.Err(err)
		}
		instr, err := instruction.NewLocalVarDef("", val)
		if err != nil {
			return errutil.Err(err)
		}
		gen.instructionBuffer = append(gen.instructionBuffer, instr)
		return nil
	default:
		gen.recurse(n)
	}
	val, err := instruction.NewGetElementPtr(i64, addr, indices)
	inst, err := instruction.NewLocalVarDef("", val)
	if err != nil {
		return errutil.Err(err)
	}
	gen.instructionBuffer = append(gen.instructionBuffer, inst)
	return nil
}

func (gen *generator) loadIndexExpr(n *ast.IndexExpr) error {
	// TODO: Fix this implementation
	// Dummy value lookup
	// i64, err := irtypes.NewInt(64)
	// if err != nil {
	// 	return errutil.Err(err)
	// }
	var nelems int
	var valtype irtypes.Type
	typ := toIrType(n.Name.Decl.Type())
	switch a := typ.(type) {
	case *irtypes.Array:
		nelems = a.Len()
		valtype = a.Elem()
	}
	val, err := instruction.NewAlloca(valtype, nelems)
	if err != nil {
		return errutil.Err(err)
	}
	//val := gen.instructionBuffer[len(gen.instructionBuffer)-1].(value.Value)
	inst, err := instruction.NewLocalVarDef("", val)
	if err != nil {
		return errutil.Err(err)
	}
	gen.instructionBuffer = append(gen.instructionBuffer, inst)
	err = gen.getElemPtr(inst, n.Index)
	if err != nil {
		return errutil.Err(err)
	}
	return nil
}

func (gen *generator) store(n ast.Expr) error {
	//valSsa := gen.ssaCounter

	//instruction.NewStore(val, addr)
	return nil
}

func (gen *generator) storeIndexExpr(n *ast.IndexExpr) error {
	// TODO: support multi dim arrays
	typ, ok := n.Name.Decl.Type().(*uctypes.Array)
	if !ok {
		errutil.Newf("indexing non array %v", n.Name.String())
	}
	arrayType := toIrType(typ)

	log.Printf("arrayType is %v", arrayType)

	// Get a pointer to the indexed element in the array
	//n.Val
	//instr:= instruction.NewGetElementPtr(elem, addr, indices)
	gen.ssaCounter++
	log.Printf("Index is stored in %v", encLocal(gen.ssaCounter))
	//gen.getElemPtr(n.Index)

	gen.ssaCounter++
	log.Printf("Create pointer to indexed position %v in array %v", encLocal(gen.ssaCounter), n.Name.String())
	return nil
}

func (gen *generator) storeIdent(n *ast.Ident) error {
	// TODO: Convert type
	switch typ := n.Decl.Type().(type) {
	case *uctypes.Basic:
		log.Printf("Storing ident %v of basic type %v", n.Name, typ.Kind)
	case *uctypes.Array:
		//if typ.Elem.(*uctypes.Basic)
		log.Printf("Storing array type %#v %T with elem %#v %T", typ, typ, typ.Elem, typ.Elem)

	}

	//log.Printf("store %v to %v", val, n.X)
	//log.Printf("Assign %v to %#v with type ", n.Y, n.X, n.X)
	return nil
}

// createBlock creates the instructions for a block statement.
func (gen *generator) createBlock(n *ast.BlockStmt) error {
	for _, blockItem := range n.Items {
		log.Printf("blockItem %#v of type %T", blockItem, blockItem)
		if err := gen.recurse(blockItem); err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// createCall creates the instructions for calling a function.
func (gen *generator) createCall(n *ast.CallExpr) error {
	var insts []instruction.Instruction
	var args []value.Value
	log.Printf("%v: create call %v with:", encLocal(gen.ssaCounter), n)
	if callType, ok := n.Name.Decl.Type().(*uctypes.Func); ok {
		// TODO: Add support for elipsis/variadic functions
		for i, arg := range n.Args {
			log.Printf("Arg %v of type %v", arg.String(), callType.Params[i].Type.String())
			if err := gen.recurse(arg); err != nil {
				return errutil.Err(err)
			}
			lastvalue, ok := gen.instructionBuffer[len(gen.instructionBuffer)-1].(value.Value)
			if !ok {
				panic("Last instruction not a value.Value")
			}
			args = append(args, lastvalue)
		}

		val, err := instruction.NewCall(toIrType(callType.Result), n.Name.String(), args)
		if err != nil {
			return errutil.Err(err)
		}

		gen.ssaCounter++
		inst, err := instruction.NewLocalVarDef("", val)
		if err != nil {
			return errutil.Err(err)
		}

		insts = append(insts, inst)
	} else {
		return errutil.Newf("invalid type assertion; expected: *types.Func, got: %T", n.Name.Decl.Type())
	}
	gen.instructionBuffer = append(gen.instructionBuffer, insts...)
	return nil
}

// createExprStmt creates the instructions for an expression statment.
func (gen *generator) createExprStmt(n *ast.ExprStmt) error {
	if err := gen.recurse(n.X); err != nil {
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
		log.Printf("create function decl %v", fn)
		return nil
	}

	log.Printf("create function def %v", fn)
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

	gen.currentFunction.SetBlocks(gen.basicBlocks)
	gen.basicBlocks = gen.basicBlocks[:0:0]

	gen.ssaCounter = 0 // TODO: check if in nested function
	gen.funcDefStack = gen.funcDefStack[:len(gen.funcDefStack)-1]
	if len(gen.funcDefStack) > 0 {
		gen.currentFunction = gen.funcDefStack[len(gen.funcDefStack)-1]
	} else {
		gen.currentFunction = nil
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
	// TODO: implement fix plix
	// if true {
	// 	newVar, err = ir.NewLocal(n.Name, ptrType)
	// } else {
	newVar, err = constant.NewPointer(ptrType, n.Name)
	// }
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
func (gen *generator) createUnaryExpr(n *ast.UnaryExpr) error {
	// TODO: implement not and minus instructions
	// Placeholder
	if err := gen.recurse(n.X); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// createVar creates the instruction for allocating place for a variable.
func (gen *generator) createVar(n *ast.VarDecl) error {
	if len(gen.funcDefStack) > 0 {
		if err := gen.createLocal(n); err != nil {
			return errutil.Err(err)
		}
	} else {
		// Global values are compile time constant, no need for ssa.
		if err := gen.createGlobal(n); err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// createWhile creates the instructions for a while loop.
func (gen *generator) createWhile(n *ast.WhileStmt) error {
	// Finnish last basic block with a branch to the basic block we
	// are about to create (with label ssaCounter+1)
	gen.ssaCounter++
	whileCond := gen.ssaCounter

	log.Printf("Clear last basic block instructions: %v", gen.instructionBuffer)
	allInsts := make([]instruction.Instruction, len(gen.instructionBuffer))
	copy(allInsts, gen.instructionBuffer)

	jmpToWhileCond, err := instruction.NewJmp(encLocal(whileCond))
	if err != nil {
		return errutil.Err(err)
	}

	bb, err := ir.NewBasicBlock(encLocal(gen.lastLabel), allInsts, jmpToWhileCond)
	gen.basicBlocks = append(gen.basicBlocks, bb)
	gen.lastLabel = gen.ssaCounter
	whileBody := gen.ssaCounter
	log.Printf("start while loop: %v", n)

	if err := gen.recurse(n.Cond); err != nil {
		return errutil.Err(err)
	}

	val, ok := gen.instructionBuffer[len(gen.instructionBuffer)-1].(value.Value)
	if !ok {
		panic("Last instruction not a valid value for branch condition")
	}
	brToWhileBody, err := instruction.NewBr(val, encLocal(whileBody), "-1")
	if err != nil {
		return errutil.Err(err)
	}
	condbb, err := ir.NewBasicBlock(encLocal(gen.lastLabel), allInsts, brToWhileBody)
	if err != nil {
		return errutil.Err(err)
	}
	gen.basicBlocks = append(gen.basicBlocks, condbb)

	gen.ssaCounter++
	gen.lastLabel = gen.ssaCounter
	gen.instructionBuffer = gen.instructionBuffer[:0]

	// Recurse over body of while loop
	log.Printf("while.Body = %#v", n.Body)

	if err := gen.recurse(n.Body); err != nil {
		return errutil.Err(err)
	}

	//_, gen.ssaCounter = endWhile(n, gen.ssaCounter)

	// End the while loop with a branch to whileCond
	bb, err = ir.NewBasicBlock(encLocal(gen.lastLabel), allInsts, jmpToWhileCond)
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

// createLocal creates the instructions for a local variable declaration.
func (gen *generator) createLocal(n *ast.VarDecl) error {
	// TODO: Implement
	var insts []instruction.Instruction
	log.Printf("create local variable %v", n)

	gen.instructionBuffer = append(gen.instructionBuffer, insts...)
	return nil
}

// createGlobal creates the instructions for a global variable declaration.
func (gen *generator) createGlobal(n *ast.VarDecl) error {
	// Ignore tentative definitions.
	if isTentativeVarDef(n) {
		log.Printf("ignoring tentative global variable definition of %v", n.Name())
		return nil
	}
	// TODO: Implement
	log.Printf("create global variable %v", n)
	var gv *ir.GlobalDecl
	gen.module.Globals = append(gen.module.Globals, gv)
	return nil
}

// isTentative reports whether the given global variable declaration is a
// tentative variable definition.
func isTentativeVarDef(n *ast.VarDecl) bool {
	ident := n.Name()
	def := ident.Decl.Name()
	return ident.Start() != def.Start()
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
		panic(errutil.Newf("Conversion failed: %#v", n))
	}
	return t
}
