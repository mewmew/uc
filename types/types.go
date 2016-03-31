package types

// TODO: Implement type checking which ensures correct uses of "void".
// Relevant sections of the uC BNF grammar have been included below.
//
//    TopLevelDecl
//       : VarDecl ";"
//       | TypeName ident "(" Formals ")" FuncBody // TypeName : "char" | "int" | "void" ;
//    ;
//
//    ScalarDecl
//       : TypeName ident // TypeName : "char" | "int" ;
//    ;
//
//    ArrayDecl
//       : TypeName ident "[" int_lit "]" // TypeName : "char" | "int" ;
//    ;
//
//    Formals
//       : TypeName   // TypeName : "void" ;
//       | FormalList
//    ;
//
//    FormalDecl
//       : ScalarDecl
//       | TypeName ident "[" "]" // TypeName : "char" | "int" ;
//    ;
