(* absyn/absyn-print.sml *)

signature ABSYN_PRINT =
  sig
    structure Absyn : ABSYN
    val program : TextIO.outstream * Absyn.program -> unit
  end (* signature ABSYN_PRINT *)

functor AbsynPrintFn(Absyn : ABSYN) : ABSYN_PRINT =
  struct

    structure Absyn = Absyn

    val say = TextIO.output
    val say1 = TextIO.output1

    fun sayIdent(os, id) = say(os, Absyn.identName id)

    fun sayInt(os, i) =
      let fun intString i =
	    let val s = Int.toString i
	    in
	      if i >= 0 then s
	      else "-" ^ String.substring(s, 1, String.size s - 1)
	    end
      in
	say(os, intString i)
      end

    fun baseTyString ty =
      case ty
	of Absyn.INTty => "int"
	 | Absyn.CHARty => "char"
	 | Absyn.VOIDty => "void"

    fun sayBaseTy(os, ty) = say(os, baseTyString ty)

    fun sayDeclarator(os, decl) =
      case decl
	of Absyn.VARdecl id => sayIdent(os, id)
	 | Absyn.ARRdecl(id,szOpt) =>
	    (say(os, "(array "); sayIdent(os, id);
	     case szOpt
	       of SOME sz => (say1(os, #" "); sayInt(os, sz))
		| NONE => ();
	     say1(os, #")"))

    fun binopString binop =
      case binop
	of Absyn.ADD => "add"
	 | Absyn.SUB => "sub"
	 | Absyn.MUL => "mul"
	 | Absyn.DIV => "div"
	 | Absyn.LT => "lt"
	 | Absyn.LE => "le"
	 | Absyn.EQ => "eq"
	 | Absyn.NE => "ne"
	 | Absyn.GE => "ge"
	 | Absyn.GT => "gt"
	 | Absyn.ANDALSO => "andalso"

    fun unopString unop =
      case unop
	of Absyn.NEG => "neg"
	 | Absyn.NOT => "not"

    fun sayConst(os, const) =
      case const
	of Absyn.INTcon i => sayInt(os, i)

    fun newLineTab(os, column) =
      let fun spaces column =
	    if column <= 0 then ()
	    else (say1(os, #" "); spaces(column-1))
	  fun tabs column =
	    if column < 8 then spaces column
	    else (say1(os, #"\t"); tabs(column-8))
      in
	say1(os, #"\n");
	tabs column
      end

    fun sayExp(os, Absyn.EXP(exp,_,_), column) =
      case exp
	of Absyn.CONST const => sayConst(os, const)
	 | Absyn.VAR id => sayIdent(os, id)
	 | Absyn.ARRAY(id,exp) =>
	    let val idExp = Absyn.EXP(Absyn.VAR id,0,0)
	    in
	      sayCall(os, "array", [idExp,exp], column)
	    end
	 | Absyn.ASSIGN(lhs,rhs) =>
	    sayCall(os, "assign", [lhs,rhs], column)
	 | Absyn.UNARY(unop, exp) =>
	    sayCall(os, unopString unop, [exp], column)
	 | Absyn.BINARY(binop,exp1,exp2) =>
	    sayCall(os, binopString binop, [exp1,exp2], column)
	 | Absyn.FCALL(id,exps) =>
	    let val idExp = Absyn.EXP(Absyn.VAR id,0,0)
	    in
	      sayCall(os, "fcall", idExp :: exps, column)
	    end

    and sayCall(os, name, args, column) =
      (say1(os, #"("); say(os, name);
       case args
	 of [] => say1(os, #")")
	  | arg::args =>
	      let val column = column + String.size name + 2
		  fun doArg arg =
		    (newLineTab(os, column); sayExp(os, arg, column))
	      in
		say1(os, #" "); sayExp(os, arg, column);
		List.app doArg args;
		say1(os, #")")
	      end)

    fun sayStmt(os, Absyn.STMT(stmt,_,_), column) =
      case stmt
	of Absyn.EMPTY =>
	    say(os, "(empty)")
	 | Absyn.EFFECT exp =>
	    sayCall(os, "effect", [exp], column)
	 | Absyn.IF(exp,thn,elsOpt) =>
	    let val column = column+2
	    in
	      say(os, "(if ");
	      sayExp(os, exp, column+2);
	      newLineTab(os, column);
	      sayStmt(os, thn, column);
	      case elsOpt
		of SOME els =>
		    (newLineTab(os, column);
		     sayStmt(os, els, column))
		 | NONE => ();
	      say1(os, #")")
	    end
	 | Absyn.WHILE(exp,stmt) =>
	    (say(os, "(while ");
	     sayExp(os, exp, column+7);
	     newLineTab(os, column+2);
	     sayStmt(os, stmt, column+2);
	     say1(os, #")"))
	 | Absyn.RETURN NONE =>
	    say(os, "(return)")
	 | Absyn.RETURN(SOME exp) =>
	    sayCall(os, "return", [exp], column)
	 | Absyn.SEQ(stmt1,stmt2) =>
	    let val column = column + 2
	    in
	      say(os, "(seq");
	      sayStmt(os, stmt1, column);
	      newLineTab(os, column);
	      sayStmt(os, stmt2, column);
	      say1(os, #")")
	    end

    fun sayVarDec(os, Absyn.VARDEC(baseTy, decl)) =
      (say1(os, #"("); sayBaseTy(os, baseTy); say1(os, #" ");
       sayDeclarator(os, decl); say1(os, #")"))

    fun saySpaceVarDec os varDec =
      (say1(os, #" "); sayVarDec(os, varDec))

    fun sayVars(os, vars) = List.app (saySpaceVarDec os) vars

    fun sayAbsFormal os (Absyn.ABSDEC(baseTy, absDecl)) =
      (say1(os, #"("); sayBaseTy(os, baseTy);
       case absDecl
	 of Absyn.EMPTYabsdecl => ()
	  | Absyn.ARRabsdecl => say(os, " array");
       say1(os, #")"))

    fun sayAbsFormals(os, absFormals) =
      List.app (sayAbsFormal os) absFormals

    fun sayDec os dec =
      case dec
	of Absyn.FUNC{name,formals,retTy,locals,body} =>
	    (say(os, "\n(function "); sayIdent(os, name); say1(os, #" "); sayBaseTy(os, retTy);
	     say(os, "\n  (formals"); sayVars(os, formals); say1(os, #")");
	     say(os, "\n  (locals"); sayVars(os, locals); say(os, ")\n  ");
	     sayStmt(os, body, 2); say(os, ")\n"))
	 | Absyn.EXTERN{name,formals,retTy} =>
	    (say(os, "\n(extern\n  (function "); sayIdent(os, name);
	     sayBaseTy(os, retTy); say(os, "\n    (formals");
	     sayVars(os, formals); say(os, ")))\n"))
	 | Absyn.GLOBAL(varDec) =>
	    (say1(os, #"\n"); sayVarDec(os, varDec); say1(os, #"\n"))

    fun program(os, Absyn.PROGRAM{decs,...}) =
      List.app (sayDec os) decs

  end (* functor AbsynPrintFn *)
