package irgen

import (
	"fmt"

	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	uctypes "github.com/mewmew/uc/types"
)

// isTentativeDef reports whether the given global variable or function
// declaration is a tentative definition.
func isTentativeDef(n ast.Decl) bool {
	ident := n.Name()
	def := ident.Decl.Name()
	return ident.Start() != def.Start()
}

// getLocal returns the LLVM IR value of the given local variable name.
func (f *Function) local(name string) value.Value {
	if v, ok := f.locals[name]; ok {
		return v
	}
	panic(fmt.Sprintf("unable to locate local variable %q", name))
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
			dbg.Printf("converting type %#v to %#v", p.Type, pt)
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
