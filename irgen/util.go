package irgen

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	uctypes "github.com/mewmew/uc/types"
)

// constZero returns the integer constant zero of the given type.
func constZero(typ irtypes.Type) constant.Constant {
	zero, err := constant.NewInt(typ, "0")
	if err != nil {
		panic(fmt.Sprintf("unable to create integer constant zero; %v", err))
	}
	return zero
}

// constOne returns the integer constant one of the given type.
func constOne(typ irtypes.Type) constant.Constant {
	one, err := constant.NewInt(typ, "1")
	if err != nil {
		panic(fmt.Sprintf("unable to create integer constant one; %v", err))
	}
	return one
}

// isTentativeDef reports whether the given global variable or function
// declaration is a tentative definition.
func isTentativeDef(n ast.Decl) bool {
	ident := n.Name()
	def := ident.Decl.Name()
	return ident.Start() != def.Start()
}

// genUnique generates a unique local variable name based on the given
// identifier.
func (f *Function) genUnique(ident *ast.Ident) string {
	name := ident.Name
	if !f.exists[name] {
		f.exists[name] = true
		return name
	}
	for i := 1; ; i++ {
		name := fmt.Sprintf("%s%d", ident.Name, i)
		if !f.exists[name] {
			f.exists[name] = true
			return name
		}
	}
	panic("unreachable")
}

// valueFromIdent returns the LLVM IR value associated with the given
// identifier. Only search for global values if f is nil.
func (m *Module) valueFromIdent(f *Function, ident *ast.Ident) value.Value {
	pos := ident.Decl.Name().Start()
	if v, ok := m.idents[pos]; ok {
		return v
	}
	if f != nil {
		if v, ok := f.idents[pos]; ok {
			return v
		}
	}
	panic(fmt.Sprintf("unable to locate value associated with identifier %q (declared at source code position %d)", ident, pos))
}

// setIdentValue maps the given global identifier to the associated value.
func (m *Module) setIdentValue(ident *ast.Ident, v value.Value) {
	pos := ident.Decl.Name().Start()
	if old, ok := m.idents[pos]; ok {
		panic(fmt.Sprintf("unable to map identifier %q to value %v; already mapped to value %v", ident, v, old))
	}
	m.idents[pos] = v
}

// setIdentValue maps the given local identifier to the associated value.
func (f *Function) setIdentValue(ident *ast.Ident, v value.Value) {
	pos := ident.Decl.Name().Start()
	if old, ok := f.idents[pos]; ok {
		panic(fmt.Sprintf("unable to map identifier %q to value %v; already mapped to value %v", ident, v, old))
	}
	f.idents[pos] = v
}

// typeOf returns the LLVM IR type of the given expression.
func (m *Module) typeOf(expr ast.Expr) irtypes.Type {
	if typ, ok := m.info.Types[expr]; ok {
		return toIrType(typ)
	}
	panic(fmt.Sprintf("unable to locate type for expression %v", expr))
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
		if ucType.Len == 0 {
			t, err = irtypes.NewPointer(elem)
		} else {
			t, err = irtypes.NewArray(elem, ucType.Len)
		}
	case *uctypes.Func:
		var params []*irtypes.Param
		variadic := false
		for _, p := range ucType.Params {
			//TODO: Add support for variadic
			if uctypes.IsVoid(p.Type) {
				break
			}
			pt := toIrType(p.Type)
			dbg.Printf("converting type %#v to %#v", p.Type, pt)
			params = append(params, irtypes.NewParam(pt, p.Name))
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
