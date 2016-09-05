// Package irgen implements a µC to LLVM IR generator.
package irgen

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/instruction"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mewkiz/pkg/errutil"
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
func (m *Module) emitGlobal(global *ir.GlobalDecl) {
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
func NewFunction(name string, sig *irtypes.Func) *Function {
	f := ir.NewFunction(name, sig)
	return &Function{Function: f, idents: make(map[int]value.Value), exists: make(map[string]bool)}
}

// startBody initializes the generation of the function body.
func (f *Function) startBody() {
	entry := f.NewBasicBlock("") // "entry"
	f.curBlock = entry
}

// endBody finalizes the generation of the function body.
func (f *Function) endBody() error {
	if block := f.curBlock; block != nil && block.Term() == nil {
		switch {
		case f.Function.Name() == "main":
			// From C11 spec $5.1.2.2.3.
			//
			// "If the return type of the main function is a type compatible with
			// int, a return from the initial call to the main function is
			// equivalent to calling the exit function with the value returned by
			// the main function as its argument; reaching the } that terminates
			// the main function returns a value of 0."
			result := f.Sig().Result()
			zero := constZero(result)
			term, err := instruction.NewRet(zero)
			if err != nil {
				panic(fmt.Sprintf("unable to create ret terminator; %v", err))
			}
			block.SetTerm(term)
		default:
			// Add void return terminator to the current basic block, if a
			// terminator is missing.
			switch result := f.Sig().Result(); {
			case irtypes.IsVoid(result):
				term, err := instruction.NewRet(nil)
				if err != nil {
					panic(fmt.Sprintf("unable to create ret instruction; %v", err))
				}
				block.SetTerm(term)
			default:
				// The semantic analysis checker guarantees that all branches of
				// non-void functions end with return statements. Therefore, if we
				// reach the current basic block doesn't have a terminator at the
				// end of the function body, it must be unreachable.
				term, err := instruction.NewUnreachable()
				if err != nil {
					panic(fmt.Sprintf("unable to create unreachable instruction; %v", err))
				}
				block.SetTerm(term)
			}
		}
	}
	f.curBlock = nil
	if err := f.AssignIDs(); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// emitInst emits to f the given unnamed value instruction.
func (f *Function) emitInst(inst instruction.ValueInst) value.Value {
	return f.curBlock.emitInst(inst)
}

// emitLocal emits to f the given named value instruction.
func (f *Function) emitLocal(ident *ast.Ident, inst instruction.ValueInst) value.Value {
	return f.curBlock.emitLocal(ident, inst)
}

// A BasicBlock represents an LLVM IR basic block generator.
type BasicBlock struct {
	// Basic block being generated.
	*ir.BasicBlock
	// Parent function of the basic block.
	parent *Function
}

// NewBasicBlock returns a new basic block generator based on the given name and
// parent function.
func (f *Function) NewBasicBlock(name string) *BasicBlock {
	block := ir.NewBasicBlock(name)
	return &BasicBlock{BasicBlock: block, parent: f}
}

// emitInst emits to b the given unnamed value instruction.
func (b *BasicBlock) emitInst(inst instruction.ValueInst) value.Value {
	v, err := instruction.NewLocalVarDef("", inst)
	if err != nil {
		panic(fmt.Sprintf("unable to create local variable definition; %v", err))
	}
	b.AppendInst(v)
	return v
}

// emitLocal emits to b the given named value instruction.
func (b *BasicBlock) emitLocal(ident *ast.Ident, inst instruction.ValueInst) value.Value {
	name := b.parent.genUnique(ident)
	v, err := instruction.NewLocalVarDef(name, inst)
	if err != nil {
		panic(fmt.Sprintf("unable to create local variable definition; %v", err))
	}
	b.AppendInst(v)
	b.parent.setIdentValue(ident, v)
	return v
}

// TODO: Consider changing the type of target from value.NamedValue to
// *ir.BasicBlock or *BasicBlock.

// emitJmp emits to b an unconditional branch terminator to the target basic
// block.
func (b *BasicBlock) emitJmp(target value.NamedValue) {
	term, err := instruction.NewJmp(target)
	if err != nil {
		panic(fmt.Sprintf("unable to create unconditional br instruction; %v", err))
	}
	b.SetTerm(term)
}

// SetTerm sets the terminator of the basic block.
func (b *BasicBlock) SetTerm(term instruction.Terminator) {
	if b.Term() != nil {
		panic(fmt.Sprintf("terminator instruction already set for basic block; old term (%v), new term (%v), basic block (%v)", term, b.Term(), b))
	}
	b.BasicBlock.SetTerm(term)
	b.parent.AppendBlock(b.BasicBlock)
}
