package irgen

import (
	"fmt"
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/instruction"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/term"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/token"
	uctypes "github.com/mewmew/uc/types"
)

// Gen generates LLVM IR based on the syntax tree of the given file.
func Gen(file *ast.File) error {
	// TODO: REMOVE log messages
	log.SetPrefix(term.BlueBold("Log:"))
	log.SetFlags(log.Lshortfile)

	var module ir.Module

	var funcDefStack []*ir.Function
	var currentFunction *ir.Function
	var basicBlocks []*ir.BasicBlock
	var instructionBuffer []instruction.Instruction
	var insts []instruction.Instruction

	instructionBuffer = make([]instruction.Instruction, 0, 10)

	// ssaCounter is counted up for anonymous assignments and basic blocks to
	// give them an unique id
	ssaCounter := 0
	lastLabel := ssaCounter
	var recurse func(ast.Node) error
	recurse = func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.File:
			for _, decl := range n.Decls {

				if err := recurse(decl); err != nil {
					return errutil.Err(err)
				}
			}
		case *ast.BlockStmt:
			for _, blockItem := range n.Items {
				log.Printf("blockItem %#v of type %t\n", blockItem, blockItem)

				if err := recurse(blockItem); err != nil {
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
			for _, prevfn := range funcDefStack {
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
			module.Funcs = append(module.Funcs, fn)

			if !astutil.IsDef(n) {
				log.Printf("create function decl %v\n", fn)
				return nil
			}

			log.Printf("create function def %v\n", fn)
			funcDefStack = append(funcDefStack, fn)
			currentFunction = fn

			// Recurse body
			if err := recurse(n.Body); err != nil {
				return errutil.Err(err)
			}
			var terminal instruction.Terminator
			terminal, ssaCounter = endFunction(fn, ssaCounter)
			allInsts := make([]instruction.Instruction, len(instructionBuffer))
			copy(allInsts, instructionBuffer)
			instructionBuffer = instructionBuffer[:0]
			basicBlocks = append(basicBlocks, ir.NewBasicBlock(toLocalVarString(ssaCounter), instructionBuffer, terminal))
			log.Print(basicBlocks)
			ssaCounter++
			funcDefStack = funcDefStack[:len(funcDefStack)-1]
			if len(funcDefStack) > 0 {
				currentFunction = funcDefStack[len(funcDefStack)-1]
			} else {
				currentFunction = nil
			}
			return nil
		case *ast.CallExpr:
			insts, ssaCounter = createCall(n, ssaCounter)
			instructionBuffer = append(instructionBuffer, insts...)
		case *ast.VarDecl:
			if uctypes.IsVoid(n.Type()) {
				return nil
			}
			if len(funcDefStack) > 0 {
				insts, ssaCounter = createLocal(n, ssaCounter)
				instructionBuffer = append(instructionBuffer, insts...)
			} else {
				// Global values are compile time constant, no need for ssa
				gv := createGlobal(n)
				module.Globals = append(module.Globals, gv)
			}
		case *ast.Ident:

			log.Println("ident")
			log.Println(n)

		case *ast.BinaryExpr:
			log.Println("bin expr")
			log.Println(n)
			if err := recurse(n.Y); err != nil {
				return errutil.Err(err)
			}
			switch n.Op {
			case token.Assign:
				log.Println(n.X)
				if x, ok := n.X.(*ast.Ident); ok {
					log.Printf("Assign %v %v %v to %#v", n.X, n.Op, n.Y, x)
					return nil
				}
			default:
				ssaCounter++
				log.Printf("Assign %v %v %v to %#v", n.X, n.Op, n.Y, toLocalVarString(ssaCounter))
			}
		case *ast.ExprStmt:
			if err := recurse(n.X); err != nil {
				return errutil.Err(err)
			}
		case *ast.WhileStmt:
			// Finnish last basic block with a branch to the basic block we
			// are about to create (with label ssaCounter+1)
			ssaCounter++
			whileLabel := ssaCounter
			allInsts := make([]instruction.Instruction, len(instructionBuffer))
			copy(allInsts, instructionBuffer)
			log.Printf("Clear last basic block instrucitons: %v\n", allInsts)
			brToWhileLabel, err := instruction.NewBr(nil, toLocalVarString(whileLabel), "")
			if err != nil {
				panic(errutil.Err(err))
			}
			//terminator, ssaCounter = createWhile(n, ssaCounter)
			basicBlocks = append(basicBlocks, ir.NewBasicBlock(toLocalVarString(lastLabel), allInsts, brToWhileLabel))

			ssaCounter++
			lastLabel = ssaCounter
			instructionBuffer = instructionBuffer[:0]

			// Recurse over body of while loop
			log.Printf("while.Body = %#v\n", n.Body)

			if err := recurse(n.Body); err != nil {
				return errutil.Err(err)
			}

			//_, ssaCounter = endWhile(n, ssaCounter)

			// End the while loop with a branch to whileLabel
			basicBlocks = append(basicBlocks, ir.NewBasicBlock(toLocalVarString(lastLabel), allInsts, brToWhileLabel))
			ssaCounter++
			lastLabel = ssaCounter
			instructionBuffer = instructionBuffer[:0]
		}
		// TODO: Implement the rest of the needed node types
		return nil
	}

	recurse(file)
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
