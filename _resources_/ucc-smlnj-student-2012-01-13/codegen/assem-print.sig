(* codegen/assem-print.sig *)

signature ASSEM_PRINT =
  sig
    structure Assem : ASSEM
    val program : TextIO.outstream * Assem.program -> unit
    val fileSuffix : string
  end (* signature ASSEM_PRINT *)
