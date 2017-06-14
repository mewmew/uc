(* rtl/rtl-print.sml *)

signature RTL_PRINT =
  sig
    structure RTL : RTL
    val program : TextIO.outstream * RTL.program -> unit
  end (* signature RTL_PRINT *)

functor RTLPrintFn(RTL : RTL) : RTL_PRINT =
  struct

    structure RTL = RTL

    val say = TextIO.output
    val say1 = TextIO.output1

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

    fun labelString(lab : RTL.label) : string = lab
    fun sayLab(os, lab) = say(os, labelString lab)

    fun relopString relop =
      case relop
	of RTL.LT => "lt"
	 | RTL.LE => "le"
	 | RTL.EQ => "eq"
	 | RTL.NE => "ne"
	 | RTL.GE => "ge"
	 | RTL.GT => "gt"

    fun tyString ty =
      case ty
	of RTL.BYTE => "_b"
	 | RTL.LONG => "_l"

    fun binopString binop =
      case binop
	of RTL.ADD => "add"
	 | RTL.SUB => "sub"
	 | RTL.MUL => "mul"
	 | RTL.DIV => "div"

    fun sayTemp(os, temp) =
      (say1(os, #"t"); sayInt(os, temp))

    fun saySpaceTemp(os, temp) =
      (say1(os, #" "); sayTemp(os, temp))

    fun sayTemps(os, temps) =
      List.app (fn temp => saySpaceTemp(os, temp)) temps

    fun sayCall(os, name, args) =
      (say1(os, #"("); say(os, name); sayTemps(os, args); say1(os, #")"))

    fun sayExpr(os, expr) =
      case expr
	of RTL.TEMP t => sayTemp(os, t)
	 | RTL.ICON i => sayInt(os, i)
	 | RTL.LABREF label => sayLab(os, label)
	 | RTL.UNARY(RTL.LOAD ty, src) =>
	    sayCall(os, "load" ^ tyString ty, [src])
	 | RTL.BINARY(binop,src1,src2) =>
	    sayCall(os, binopString binop, [src1,src2])

    fun sayInsn os insn =
      case insn
	of RTL.LABDEF label =>
	    (sayLab(os, label); say(os, ":\n"))
	 | RTL.JUMP label =>
	    (say(os, "\t(goto "); sayLab(os, label); say(os, ")\n"))
	 | RTL.CJUMP(relop,src1,src2,label) =>
	    (say(os, "\t(if ");
	     sayCall(os, relopString relop, [src1,src2]);
	     say(os, " (goto "); sayLab(os, label); say(os, "))\n"))
	 | RTL.STORE(ty,dst,src) =>
	    (say1(os, #"\t"); sayCall(os, "store" ^ tyString ty, [dst,src]); say1(os, #"\n"))
	 | RTL.EVAL(dst,expr) =>
	    (say(os, "\t(set "); sayTemp(os, dst); say1(os, #" ");
	     sayExpr(os, expr); say(os, ")\n"))
	 | RTL.CALL(SOME dst,label,args) =>
	    (say1(os, #"\t"); say(os, "(set "); sayTemp(os, dst); say1(os, #" ");
	     sayCall(os, "call "^label, args); say(os, ")\n"))
	 | RTL.CALL(NONE,label,args) =>
	    (say1(os, #"\t"); sayCall(os, "call "^label, args); say1(os, #"\n"))

    fun sayDec os dec =
      case dec
	of RTL.PROC{label,formals,locals,frameSize,insns} =>
	    (say(os, "\n(procedure "); sayLab(os, label);
	     say(os, "\n\t(formals"); sayTemps(os, formals); say1(os, #")");
	     say(os, "\n\t(locals"); sayTemps(os, locals); say1(os, #")");
	     say(os, "\n\t(frameSize "); sayInt(os, frameSize); say(os, ")\n");
	     List.app (sayInsn os) insns; say(os, ")\n"))
	 | RTL.DATA{label,size} =>
	    (say(os, "\n(data "); sayLab(os, label);
	     say1(os, #" "); sayInt(os,size); say(os, ")\n"))

    fun program(os, RTL.PROGRAM decs) =
      List.app (sayDec os) decs

  end (* functor RTLPrintFn *)
