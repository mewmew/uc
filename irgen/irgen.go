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
	"github.com/mewmew/uc/types"
)

func Gen(file *ast.File) error {
	// TODO: REMOVE
	log.SetPrefix(term.BlueBold("Log:"))
	log.SetFlags(log.Lshortfile)

	var module ir.Module

	var functions []*ir.Function
	var currentFunction *ir.Function
	var basicBlocks []*ir.BasicBlock
	var instructionBuffer []instruction.Instruction
	var insts []instruction.Instruction

	instructionBuffer = make([]instruction.Instruction, 10)

	// ssaCounter is counted up for anonymous assignments and basic blocks to
	// give them an unique id
	ssaCounter := 0

	genBefore := func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.FuncDecl:
			var fn *ir.Function
			fn, ssaCounter = createFunction(n, ssaCounter)
			module.Funcs = append(module.Funcs, fn)
			if astutil.IsDef(n) {
				functions = append(functions, fn)
				currentFunction = fn
			}
		case *ast.CallExpr:
			insts, ssaCounter = createCall(n, ssaCounter)
			instructionBuffer = append(instructionBuffer, insts...)
		case *ast.VarDecl:
			if types.IsVoid(n.Type()) {
				return nil
			}
			if len(functions) > 0 {
				insts, ssaCounter = createLocal(n, ssaCounter)
				instructionBuffer = append(instructionBuffer, insts...)
			} else {
				// Global values are compile time constant, no need for ssa
				gv := createGlobal(n)
				module.Globals = append(module.Globals, gv)
			}
		case *ast.WhileStmt:
			// TODO: Create branch and 2 new basic blocks
			allInsts := make([]instruction.Instruction, len(instructionBuffer))
			copy(allInsts, instructionBuffer)
			log.Printf("All basic block instrucitons: %v\n", allInsts)
			branch, ssaCounter := createWhile(n, ssaCounter)
			basicBlocks = append(basicBlocks, ir.NewBasicBlock(toLocalVarString(ssaCounter), instructionBuffer, branch))
			ssaCounter++
			instructionBuffer = instructionBuffer[:0]
		}
		// TODO: Implement the rest of the needed node types
		return nil
	}

	genAfter := func(n ast.Node) error {

		switch n := n.(type) {
		case *ast.WhileStmt:
			_, ssaCounter = endWhile(n, ssaCounter)
		case *ast.FuncDecl:
			if astutil.IsDef(n) {
				terminal, ssaCounter := endFunction(n, ssaCounter)
				allInsts := make([]instruction.Instruction, len(instructionBuffer))
				copy(allInsts, instructionBuffer)
				basicBlocks = append(basicBlocks, ir.NewBasicBlock(toLocalVarString(ssaCounter), instructionBuffer, terminal))
				log.Print(basicBlocks)
				ssaCounter++
				functions = functions[:len(functions)-1]
				if len(functions)-1 > 0 {
					currentFunction = functions[len(functions)-1]
				} else {
					currentFunction = nil
				}
			}

		}
		return nil
	}

	// Walk the AST of the given file to generate IR.
	if err := astutil.WalkBeforeAfter(file, genBefore, genAfter); err != nil {
		return errutil.Err(err)
	}

	return nil
}

func createFunction(fn *ast.FuncDecl, ssa int) (*ir.Function, int) {
	// TODO: Implement
	log.Printf("create function decl %v\n", fn)
	return nil, ssa
}

func endFunction(fn *ast.FuncDecl, ssa int) (instruction.Terminator, int) {
	// TODO: Implement
	log.Printf("end function decl %v\n", fn)
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
	log.Printf("%v: create local variable %v\n", toLocalVarString(ssa), lv)
	ssa++
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
