(* mips/mips-assem-print.sml *)

functor MIPSAssemPrintFn(structure MIPS : MIPS
			   ) : ASSEM_PRINT =
  struct

    structure Assem = MIPS

    fun program(os, p) = MIPS.print_program(p, os)

    val fileSuffix = ".s"

  end (* functor MIPSAssemPrintFn *)
