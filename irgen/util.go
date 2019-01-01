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

// implicitConversion implicitly converts the value of the smallest type to the
// largest type of x and y, emitting code to f. The new values of x and y are
// returned.
func (m *Module) implicitConversion(f *Func, x, y value.Value) (value.Value, value.Value) {
	// Implicit conversion.
	switch {
	case isLarger(x.Type(), y.Type()):
		y = m.convert(f, y, x.Type())
	case isLarger(y.Type(), x.Type()):
		x = m.convert(f, x, y.Type())
	}
	return x, y
}

// convert converts the given value to the specified type, emitting code to f.
// No conversion is made, if v is already of the correct type.
func (m *Module) convert(f *Func, v value.Value, to irtypes.Type) value.Value {
	// Early return if v is already of the correct type.
	from := v.Type()
	if irtypes.Equal(from, to) {
		return v
	}
	fromType, ok := from.(*irtypes.IntType)
	if !ok {
		panic(fmt.Sprintf("support for converting from type %T not yet implemented", from))
	}
	toType, ok := to.(*irtypes.IntType)
	if !ok {
		panic(fmt.Sprintf("support for converting to type %T not yet implemented", to))
	}

	// Convert constant values.
	switch v := v.(type) {
	case *constant.Int:
		c, err := constant.NewIntFromString(toType, v.Ident())
		if err != nil {
			panic(fmt.Errorf("unable to parse integer literal %q; %v", v.Ident(), err))
		}
		return c
	case *constant.Float:
		panic(fmt.Sprintf("support for converting type %T not yet implemented", v))
	}

	// TODO: Add proper support for converting signed and unsigned values, using
	// sext and zext, respectively.

	// Convert unsigned values.
	if fromType.Equal(irtypes.I1) {
		// Zero extend boolean values.
		return f.curBlock.NewZExt(v, toType)
	}

	// Convert signed values.
	if toType.BitSize > fromType.BitSize {
		// Sign extend.
		return f.curBlock.NewSExt(v, toType)
	}
	// Truncate.
	return f.curBlock.NewTrunc(v, toType)
}

// isLarger reports whether t has higher precision than u.
func isLarger(t, u irtypes.Type) bool {
	// A Sizer is a type with a size in number of bits.
	type Sizer interface {
		// Size returns the size of t in number of bits.
		Size() int
	}
	if t, ok := t.(Sizer); ok {
		if u, ok := u.(Sizer); ok {
			return t.Size() > u.Size()
		}
	}
	return false
}

// isRef reports whether the given type is a reference type; e.g. pointer or
// array.
func isRef(typ irtypes.Type) bool {
	switch typ.(type) {
	case *irtypes.ArrayType:
		return true
	case *irtypes.PointerType:
		return true
	default:
		return false
	}
}

// constZero returns the integer constant zero of the given type.
func constZero(typ irtypes.Type) constant.Constant {
	intType, ok := typ.(*irtypes.IntType)
	if !ok {
		panic(fmt.Errorf("invalid integer literal type; expected *types.IntType, got %T", typ))
	}
	return constant.NewInt(intType, 0)
}

// constOne returns the integer constant one of the given type.
func constOne(typ irtypes.Type) constant.Constant {
	intType, ok := typ.(*irtypes.IntType)
	if !ok {
		panic(fmt.Errorf("invalid integer literal type; expected *types.IntType, got %T", typ))
	}
	return constant.NewInt(intType, 1)
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
func (f *Func) genUnique(ident *ast.Ident) string {
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
}

// isGlobal reports whether the given identifier is a global definition.
func (m *Module) isGlobal(ident *ast.Ident) bool {
	pos := ident.Decl.Name().Start()
	_, exists := m.idents[pos]
	return exists
}

// valueFromIdent returns the LLVM IR value associated with the given
// identifier. Only search for global values if f is nil.
func (m *Module) valueFromIdent(f *Func, ident *ast.Ident) value.Value {
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
func (f *Func) setIdentValue(ident *ast.Ident, v value.Value) {
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

// toIrType converts the given uC type to the corresponding LLVM IR type.
func toIrType(n uctypes.Type) irtypes.Type {
	//TODO: implement, placeholder implementation
	var t irtypes.Type
	var err error
	switch ucType := n.(type) {
	case *uctypes.Basic:
		switch ucType.Kind {
		case uctypes.Int:
			//TODO: Get int width from compile env
			t = irtypes.NewInt(32)
		case uctypes.Char:
			t = irtypes.NewInt(8)
		case uctypes.Void:
			t = irtypes.Void
		}
	case *uctypes.Array:
		elem := toIrType(ucType.Elem)
		if ucType.Len == 0 {
			t = irtypes.NewPointer(elem)
		} else {
			t = irtypes.NewArray(uint64(ucType.Len), elem)
		}
	case *uctypes.Func:
		var params []irtypes.Type
		variadic := false
		for _, p := range ucType.Params {
			//TODO: Add support for variadic
			if uctypes.IsVoid(p.Type) {
				break
			}
			pt := toIrType(p.Type)
			dbg.Printf("converting type %#v to %#v", p.Type, pt)
			params = append(params, pt)
		}
		result := toIrType(ucType.Result)
		typ := irtypes.NewFunc(result, params...)
		typ.Variadic = variadic
		t = typ
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
