(* mips/mips-codegen.sml *)

functor MIPSCodegenFn(structure RTL : RTL
		      structure MIPS : MIPS
		      (* YOUR HELPER MODULES HERE *)
			) : CODEGEN =
  struct

    structure RTL = RTL
    structure Assem = MIPS

    fun program p = []	(* XXX: REPLACE WITH YOUR CODE *)

  end (* functor MIPSCodegenFn *)
