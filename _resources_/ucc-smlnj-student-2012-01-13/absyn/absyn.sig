(* absyn/absyn.sig *)

signature ABSYN =
  sig

    structure Source : SOURCE

    datatype ident
      = IDENT of string * int * int

    (* TYPES & DECLARATORS *)

    datatype baseTy
      = INTty
      | CHARty
      | VOIDty

    datatype declarator
      = VARdecl of ident
      | ARRdecl of ident * int option

    (* EXPRESSIONS *)

    datatype binop
      = ADD | SUB | MUL | DIV
      | LT | LE | EQ | NE | GE | GT
      | ANDALSO 

    datatype unop
      = NEG | NOT

    datatype const
      = INTcon of int

    datatype exp
      = EXP of exp' * int * int
    and exp'
      = CONST of const
      | VAR of ident
      | ARRAY of ident * exp
      | ASSIGN of exp * exp
      | UNARY of unop * exp
      | BINARY of binop * exp * exp
      | FCALL of ident * exp list

    (* STATEMENTS *)

    datatype stmt
      = STMT of stmt' * int * int
    and stmt'
      = EMPTY
      | EFFECT of exp
      | IF of exp * stmt * stmt option
      | WHILE of exp * stmt
      | RETURN of exp option
      | SEQ of stmt * stmt

    (* DECLARATIONS *)

    datatype varDec
      = VARDEC of baseTy * declarator

    datatype absDecl
      = EMPTYabsdecl
      | ARRabsdecl

    datatype absFormal
      = ABSDEC of baseTy * absDecl

    datatype topDec
      = FUNC of {name: ident,
		 formals: varDec list,
		 retTy: baseTy,
		 locals: varDec list,
		 body: stmt}
      | EXTERN of {name: ident,
		   formals: varDec list,
		   retTy: baseTy}
      | GLOBAL of varDec

    datatype program
      = PROGRAM of {decs: topDec list,
		    source: Source.source}

    structure IdentDict : ORD_DICT where type Key.ord_key = ident
    val makeIdent	: string * int * int -> ident
    val identName	: ident -> string
    val identEqual	: ident * ident -> bool

  end (* signature ABSYN *)
