package irgen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/instruction"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/sem"
)

// A Module represents an LLVM IR module generator.
type Module struct {
	// Module being generated.
	*ir.Module
	// info holds semantic information about the program from the type-checker.
	info *sem.Info
}

// NewModule returns a new module generator.
func NewModule(info *sem.Info) *Module {
	m := ir.NewModule()
	return &Module{Module: m, info: info}
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
}

// NewFunction returns a new function generator based on the given function name
// and signature.
//
// The caller is responsible for initializing basic blocks.
func NewFunction(name string, sig *irtypes.Func) *Function {
	f := ir.NewFunction(name, sig)
	return &Function{Function: f}
}

// startBody initializes the generation of the function body.
func (f *Function) startBody() {
	entry := NewBasicBlock("") // "entry"
	f.curBlock = entry
}

// endBody finalizes the generation of the function body.
func (f *Function) endBody() error {
	f.curBlock = nil
	if err := f.AssignIDs(); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// emitInst emits to f the given instruction.
func (f *Function) emitInst(inst instruction.ValueInst) value.Value {
	return f.curBlock.emitInst(inst)
}

// A BasicBlock represents an LLVM IR basic block generator.
type BasicBlock struct {
	// Basic block being generated.
	*ir.BasicBlock
}

// NewBasicBlock returns a new basic block generator.
func NewBasicBlock(name string) *BasicBlock {
	block := ir.NewBasicBlock(name)
	return &BasicBlock{BasicBlock: block}
}

// emitInst emits to b the given instruction.
func (b *BasicBlock) emitInst(inst instruction.ValueInst) value.Value {
	def := instruction.NewLocalVarDef("", inst)
	b.AppendInst(def)
	return def
}
