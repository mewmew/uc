(* codegen/codegen.sig *)

signature CODEGEN =
  sig
    structure RTL : RTL
    structure Assem : ASSEM
    val program : RTL.program -> Assem.program
  end (* signature CODEGEN *)
