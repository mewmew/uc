// Package irgen implements a µC to LLVM IR generator.
package irgen

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/llir/llvm/ir"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mewkiz/pkg/term"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/sem"
)

// TODO: Remove debug output.

// dbg is a logger which prefixes debug messages with "irgen:".
var dbg = log.New(ioutil.Discard, term.WhiteBold("irgen:"), log.Lshortfile)

// A Module represents an LLVM IR module generator.
type Module struct {
	// Module being generated.
	*ir.Module
	// info holds semantic information about the program from the type-checker.
	info *sem.Info
	// Maps from identifier source code position to the associated value.
	idents map[int]value.Value
}

// NewModule returns a new module generator.
func NewModule(info *sem.Info) *Module {
	m := ir.NewModule()
	return &Module{Module: m, info: info, idents: make(map[int]value.Value)}
}

// emitFunc emits to m the given function.
func (m *Module) emitFunc(f *Function) {
	m.Funcs = append(m.Funcs, f.Function)
}

// emitGlobal emits to m the given global variable declaration.
func (m *Module) emitGlobal(global *ir.Global) {
	m.Globals = append(m.Globals, global)
}

// A Function represents an LLVM IR function generator.
type Function struct {
	// Function being generated.
	*ir.Function
	// Current basic block being generated.
	curBlock *BasicBlock
	// Maps from identifier source code position to the associated value.
	idents map[int]value.Value
	// Map of existing local variable names.
	exists map[string]bool
}

// NewFunction returns a new function generator based on the given function name
// and signature.
//
// The caller is responsible for initializing basic blocks.
func NewFunction(name string, sig *irtypes.FuncType) *Function {
	var params []*ir.Param
	for _, paramType := range sig.Params {
		param := ir.NewParam("", paramType)
		params = append(params, param)
	}
	f := ir.NewFunction(name, sig.RetType, params...)
	return &Function{Function: f, idents: make(map[int]value.Value), exists: make(map[string]bool)}
}

// startBody initializes the generation of the function body.
func (f *Function) startBody() {
	entry := f.NewBlock("") // "entry"
	f.curBlock = entry
}

// endBody finalizes the generation of the function body.
func (f *Function) endBody() error {
	if block := f.curBlock; block != nil && block.Term == nil {
		switch {
		case f.Function.GlobalName == "main":
			// From C11 spec $5.1.2.2.3.
			//
			// "If the return type of the main function is a type compatible with
			// int, a return from the initial call to the main function is
			// equivalent to calling the exit function with the value returned by
			// the main function as its argument; reaching the } that terminates
			// the main function returns a value of 0."
			result := f.Sig.RetType
			zero := constZero(result)
			termRet := ir.NewRet(zero)
			block.SetTerm(termRet)
		default:
			// Add void return terminator to the current basic block, if a
			// terminator is missing.
			switch result := f.Sig.RetType; {
			case result.Equal(irtypes.Void):
				termRet := ir.NewRet(nil)
				block.SetTerm(termRet)
			default:
				// The semantic analysis checker guarantees that all branches of
				// non-void functions end with return statements. Therefore, if we
				// reach the current basic block doesn't have a terminator at the
				// end of the function body, it must be unreachable.
				termUnreachable := ir.NewUnreachable()
				block.SetTerm(termUnreachable)
			}
		}
	}
	f.curBlock = nil
	return nil
}

// emitLocal emits to f the given named value instruction.
func (f *Function) emitLocal(ident *ast.Ident, inst valueInst) value.Value {
	return f.curBlock.emitLocal(ident, inst)
}

// A BasicBlock represents an LLVM IR basic block generator.
type BasicBlock struct {
	// Basic block being generated.
	*ir.BasicBlock
	// Parent function of the basic block.
	parent *Function
}

// NewBlock returns a new basic block generator based on the given name and
// parent function.
func (f *Function) NewBlock(name string) *BasicBlock {
	block := ir.NewBlock(name)
	return &BasicBlock{BasicBlock: block, parent: f}
}

// valueInst represents an instruction producing a value.
type valueInst interface {
	ir.Instruction
	value.Named
}

// emitLocal emits to b the given named value instruction.
func (b *BasicBlock) emitLocal(ident *ast.Ident, inst valueInst) value.Value {
	name := b.parent.genUnique(ident)
	inst.SetName(name)
	b.parent.setIdentValue(ident, inst)
	return inst
}

// SetTerm sets the terminator of the basic block.
func (b *BasicBlock) SetTerm(term ir.Terminator) {
	if b.Term != nil {
		panic(fmt.Sprintf("terminator instruction already set for basic block; old term (%v), new term (%v), basic block (%v)", term, b.Term, b))
	}
	b.BasicBlock.Term = term
	b.parent.Blocks = append(b.parent.Blocks, b.BasicBlock)
}
