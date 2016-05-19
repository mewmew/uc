package irgen

import (
	"fmt"
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/instruction"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/term"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/types"
)

var SsaCounter int

func Gen(file *ast.File) error {
	// TODO: REMOVE
	log.SetPrefix(term.BlueBold("Log:"))
	log.SetFlags(log.Lshortfile)

	var module ir.Module

	var functions []*ir.Function
	var currentFunction *ir.Function
	var basicBlocks []*ir.BasicBlock
	var instructionBuffer []instruction.Instruction

	instructionBuffer = make([]instruction.Instruction, 10)
	SsaCounter := 0

	genBefore := func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.FuncDecl:
			fn := createFunction(n)
			module.Funcs = append(module.Funcs, fn)
			if astutil.IsDef(n) {
				functions = append(functions, fn)
				currentFunction = fn
			}
		case *ast.CallExpr:
			insts := createCall(n)
			instructionBuffer = append(instructionBuffer, insts...)
		case *ast.VarDecl:
			if types.IsVoid(n.Type()) {
				return nil
			}
			if len(functions) > 0 {
				insts := createLocal(n)
				instructionBuffer = append(instructionBuffer, insts...)
			} else {
				gv := createGlobal(n)
				module.Globals = append(module.Globals, gv)
			}
			SsaCounter += 1
		case *ast.WhileStmt:
			_ = basicBlocks
			// TODO: Create branch and 2 new basic blocks
			allInsts := make([]instruction.Instruction, len(instructionBuffer))
			copy(allInsts, instructionBuffer)
			log.Printf("All basic block instrucitons: %v\n", allInsts)
			branch := createWhile(n)
			basicBlocks = append(basicBlocks, ir.NewBasicBlock(toLocalVarString(SsaCounter), instructionBuffer, branch))
			instructionBuffer = instructionBuffer[0:0]
		}
		// TODO: Implement the rest of the needed node types
		return nil
	}

	genAfter := func(n ast.Node) error {

		switch n := n.(type) {
		case *ast.WhileStmt:
			endWhile(n)
		case *ast.FuncDecl:
			if astutil.IsDef(n) {
				insts := endFunction(n)
				instructionBuffer = append(instructionBuffer, insts...)
				l := len(functions)
				functions = functions[:l-1]
				l = len(functions)
				if l-1 > 0 {
					currentFunction = functions[l-1]
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

	_ = ir.BasicBlock{}
	return nil
}

func createFunction(fn *ast.FuncDecl) *ir.Function {
	// TODO: Implement
	log.Println("create function decl", fn)
	return nil
}

func endFunction(fn *ast.FuncDecl) []instruction.Instruction {
	// TODO: Implement
	log.Println("end function decl", fn)
	return nil
}

func createCall(call *ast.CallExpr) []instruction.Instruction {
	// TODO: Implement
	log.Println("create call", call)
	return nil
}

func createLocal(lv *ast.VarDecl) []instruction.Instruction {
	// TODO: Implement
	log.Println("create local variable", lv)
	SsaCounter += 1
	return nil
}

func createGlobal(gv *ast.VarDecl) *ir.GlobalDecl {
	// TODO: Implement
	log.Println("create global variable", gv)
	SsaCounter += 1
	return nil
}

func createWhile(gv *ast.WhileStmt) instruction.Terminator {
	// TODO: Implement
	log.Println("start while loop", gv)
	return nil
}

func endWhile(gv *ast.WhileStmt) []instruction.Instruction {
	// TODO: Implement
	log.Println("end while loop", gv)
	return nil
}

func toLocalVarString(ssa int) string {
	return fmt.Sprintf("%%%v", ssa)
}
