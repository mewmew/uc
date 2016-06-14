package irgen

import (
	"fmt"
	"log"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/instruction"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/term"
	"github.com/mewmew/uc/sem"
)

// dbg is a logger which prefixes debug messages with "irgen:".
var dbg *log.Logger

// A Module represents an LLVM IR module generator.
type Module struct {
	// Module being generated.
	*ir.Module
	// info holds semantic information about the program from the type-checker.
	info *sem.Info
}

// NewModule returns a new module generator.
func NewModule(info *sem.Info) *Module {
	// TODO: Remove debug output.
	dbg = log.New(os.Stderr, term.WhiteBold("irgen:"), log.Lshortfile)
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
	// Local variables.
	local map[string]value.Value
}

// NewFunction returns a new function generator based on the given function name
// and signature.
//
// The caller is responsible for initializing basic blocks.
func NewFunction(name string, sig *irtypes.Func) *Function {
	f := ir.NewFunction(name, sig)
	return &Function{Function: f, local: make(map[string]value.Value)}
}

// startBody initializes the generation of the function body.
func (f *Function) startBody() {
	entry := f.NewBasicBlock("") // "entry"
	f.curBlock = entry
}

// endBody finalizes the generation of the function body.
func (f *Function) endBody() error {
	if block := f.curBlock; block != nil && block.Term() == nil {
		// Add void return terminator to the current basic block, if a terminator
		// is missing.
		if result := f.Type().Result(); !irtypes.IsVoid(result) {
			panic(fmt.Sprintf("unable to finalize current basic block of function body; expected void return since terminator was missing, got %v", result))
		}
		term, err := instruction.NewRet(irtypes.NewVoid(), nil)
		if err != nil {
			panic(fmt.Sprintf("unable to create ret instruction; %v", err))
		}
		block.SetTerm(term)
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
func (f *Function) emitLocal(name string, inst instruction.ValueInst) value.Value {
	return f.curBlock.emitLocal(name, inst)
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
	f.AppendBlock(block)
	return &BasicBlock{BasicBlock: block, parent: f}
}

// emitInst emits to b the given unnamed value instruction.
func (b *BasicBlock) emitInst(inst instruction.ValueInst) value.Value {
	return b.emitLocal("", inst)
}

// emitLocal emits to b the given named value instruction.
func (b *BasicBlock) emitLocal(name string, inst instruction.ValueInst) value.Value {
	def := instruction.NewLocalVarDef(name, inst)
	b.AppendInst(def)
	b.parent.local[name] = def
	return def
}
