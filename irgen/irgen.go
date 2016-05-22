package irgen

import (
	"fmt"
	"log"

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
	// lastLabel holds the ssa of the last created basic block
	recurse func(ast.Node) error
}

// Gen generates LLVM IR based on the syntax tree of the given file.
func Gen(file *ast.File) error {
	// TODO: REMOVE log messages
	log.SetPrefix(term.BlueBold("Log:"))
	log.SetFlags(log.Lshortfile)
	var gen generator
	// The usual size of a basic block seems to be less than 10 instructions.
	gen.instructionBuffer = make([]instruction.Instruction, 0, 10)

	var insts []instruction.Instruction

	gen.ssaCounter = 0
	gen.lastLabel = gen.ssaCounter
	gen.recurse = func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.File:
			for _, decl := range n.Decls {

				if err := gen.recurse(decl); err != nil {
					return errutil.Err(err)
				}
			}
		case *ast.BlockStmt:
			for _, blockItem := range n.Items {
				log.Printf("blockItem %#v of type %T\n", blockItem, blockItem)

				if err := gen.recurse(blockItem); err != nil {
					return errutil.Err(err)
				}
			}
		case *ast.FuncDecl:
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
			var terminal instruction.Terminator
			terminal, gen.ssaCounter = endFunction(fn, gen.ssaCounter)
			allInsts := make([]instruction.Instruction, len(gen.instructionBuffer))
			copy(allInsts, gen.instructionBuffer)
			gen.instructionBuffer = gen.instructionBuffer[:0]
			bb, err := ir.NewBasicBlock(toLocalVarString(gen.lastLabel), allInsts, terminal)
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
		case *ast.CallExpr:
			insts, gen.ssaCounter = createCall(n, gen.ssaCounter)
			gen.instructionBuffer = append(gen.instructionBuffer, insts...)
		case *ast.VarDecl:
			if uctypes.IsVoid(n.Type()) {
				return nil
			}
			if len(gen.funcDefStack) > 0 {
				insts, gen.ssaCounter = createLocal(n, gen.ssaCounter)
				gen.instructionBuffer = append(gen.instructionBuffer, insts...)
			} else {
				// Global values are compile time constant, no need for ssa
				// NOTE: Global variables may still be unnamed, but they would have
				// a different counter; e.g. @0 and %0 may co-exist. We may split
				// ssaCounter into a local and a global counter, along the lines of:
				//
				//    localCounter := 0
				//    globalCounter := 0
				gv := createGlobal(n)
				gen.module.Globals = append(gen.module.Globals, gv)
			}
		case *ast.Ident:
			gen.ssaCounter++
			log.Printf("loading ident %v into %v", n.Name, toLocalVarString(gen.ssaCounter))
			typ := toIrType(n.Decl.Type())
			var newVar value.Value
			var err error
			// TODO: add data fields or deduce if n is local or global

			if true {
				newVar, err = ir.NewLocal(n.Name, typ)
			} else {
				var ptrType *irtypes.Pointer
				ptrType, err = irtypes.NewPointer(typ)
				if err != nil {
					return errutil.Err(err)
				}
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
		case *ast.BinaryExpr:
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
			log.Printf("Assign %v %v %v to %#v", n.X, n.Op, n.Y, toLocalVarString(gen.ssaCounter))

		case *ast.ExprStmt:
			if err := gen.recurse(n.X); err != nil {
				return errutil.Err(err)
			}
		case *ast.WhileStmt:
			// Finnish last basic block with a branch to the basic block we
			// are about to create (with label ssaCounter+1)
			gen.ssaCounter++
			whileLabel := gen.ssaCounter
			allInsts := make([]instruction.Instruction, len(gen.instructionBuffer))
			copy(allInsts, gen.instructionBuffer)
			log.Printf("Clear last basic block instrucitons: %v\n", allInsts)
			brToWhileLabel, err := instruction.NewBr(nil, toLocalVarString(whileLabel), "")
			if err != nil {
				return errutil.Err(err)
			}
			//terminator, gen.ssaCounter = createWhile(n, gen.ssaCounter)
			bb, err := ir.NewBasicBlock(toLocalVarString(gen.lastLabel), allInsts, brToWhileLabel)
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
			bb, err = ir.NewBasicBlock(toLocalVarString(gen.lastLabel), allInsts, brToWhileLabel)
			if err != nil {
				return errutil.Err(err)
			}
			gen.basicBlocks = append(gen.basicBlocks, bb)
			gen.ssaCounter++
			gen.lastLabel = gen.ssaCounter
			gen.instructionBuffer = gen.instructionBuffer[:0]
		}
		// TODO: Implement the rest of the needed node types
		return nil
	}

	gen.recurse(file)
	return nil
}

func createFunction(fn *ast.FuncDecl, ssa int) (*ir.Function, int) {
	// TODO: Implement
	//instr := ir.NewFunction(fn.Name(), )
	log.Printf("create function decl %v\n", fn)
	return nil, ssa
}

func endFunction(fn *ir.Function, ssa int) (instruction.Terminator, int) {
	// TODO: Implement
	log.Printf("end function def %v\n", fn)
	ret, err := instruction.NewRet(irtypes.NewVoid(), nil)
	if err != nil {
		log.Panic(errutil.New(err.Error()))
	}
	// TODO: how to return 0 from main without return stmt
	return ret, ssa
}

func createCall(call *ast.CallExpr, ssa int) ([]instruction.Instruction, int) {
	// TODO: Implement
	log.Printf("%v: create call %v\n", toLocalVarString(ssa), call)
	ssa++
	return nil, ssa
}

func createLocal(lv *ast.VarDecl, ssa int) ([]instruction.Instruction, int) {
	// TODO: Implement
	log.Printf("create local variable %v\n", lv)
	return nil, ssa
}

func createGlobal(gv *ast.VarDecl) *ir.GlobalDecl {
	// TODO: Implement
	log.Printf("create global variable %v\n", gv)
	return nil
}

func createWhile(ws *ast.WhileStmt, ssa int) (instruction.Terminator, int) {
	// TODO: Implement

	log.Printf("start while loop %v\n", ws)
	return nil, ssa
}

func endWhile(gv *ast.WhileStmt, ssa int) ([]instruction.Instruction, int) {
	// TODO: Implement
	log.Printf("end while loop %v\n", gv)
	return nil, ssa
}

// NOTE: Replace with asm.EncLocal.
func toLocalVarString(ssa int) string {
	return fmt.Sprintf("%%%v", ssa)
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
