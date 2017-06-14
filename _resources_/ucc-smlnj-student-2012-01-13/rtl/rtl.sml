(* rtl/rtl.sml *)

structure RTL : RTL =
  struct

    datatype ty	(* for memory load/store operations *)
      = BYTE	(* char *)
      | LONG	(* int or pointer *)

    type temp = int	(* pseudo register *)
    type label = string

    val RV = 0	(* temp for return value *)
    val FP = 1	(* frame pointer temp, for storage of local arrays *)

    datatype relop
      = LT | LE | EQ | NE | GE | GT

    datatype unop
      = LOAD of ty	(* load value of type 'ty' from address in operand *)

    datatype binop
      = ADD | SUB | MUL | DIV

    datatype expr
      (* moving values into the result temp *)
      = TEMP of temp
      | ICON of int
      | LABREF of label
      (* unary operators *)
      | UNARY of unop * temp
      (* binary operators *)
      | BINARY of binop * temp * temp

    datatype insn
      (* control flow *)
      = LABDEF of label
      | JUMP of label
      | CJUMP of relop * temp * temp * label
      (* stores to memory *)
      | STORE of ty * temp * temp
      (* simple expression evaluation *)
      | EVAL of temp * expr
      (* function calls: could be expr but result temp is optional *)
      | CALL of temp option * label * temp list

    datatype dec
      = PROC of {label: label, formals: temp list, locals: temp list,
		 frameSize: int, insns: insn list}
      | DATA of {label: label, size: int}

    datatype program
      = PROGRAM of dec list

    fun tick counter =
      let val i = !counter + 1
	  val _ = counter := i
      in
	i
      end

    val labCounter = ref 99
    fun newLabel() = "L" ^ Int.toString(tick labCounter)

    val tempCounter = ref 1
    fun newTemp() = tick tempCounter

    fun sizeof(BYTE) = 1	(* per definition *)
      | sizeof(LONG) = 4	(* ok for 32-bitters *)

  end (* structure RTL *)
