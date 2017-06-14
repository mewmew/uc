(* codegen/insns.sig *)

signature INSNS =
  sig

    (* This can be used by a backend to interface
       with a graph-colouring register allocator. *)

    type move
    type oper
    type jump

    type temp = int
    type label = string
    val moveDefUse : move -> {def: temp, use: temp}
    val operDefUse : oper -> {def: temp list, use: temp list}
    val jumpDefUse : jump -> {def: temp list, use: temp list}
    val jumpTargets: jump -> {labels: label list, fallThrough: bool}

    datatype insn
      = LABEL of label
      | MOVE of move
      | OPER of oper
      | JUMP of jump

  end (* signature INSNS *)
