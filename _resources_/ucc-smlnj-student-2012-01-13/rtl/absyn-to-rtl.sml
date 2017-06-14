(* rtl/absyn-to-rtl.sml *)

signature ABSYN_TO_RTL =
  sig
    structure Absyn	: ABSYN
    structure RTL	: RTL
    val program		: Absyn.program -> RTL.program
  end (* signature ABSYN_TO_RTL *)

functor AbsynToRTLFn(structure Absyn : ABSYN
		     structure RTL : RTL
		       ) : ABSYN_TO_RTL =
  struct

    structure Absyn = Absyn
    structure RTL = RTL

    structure IdentDict = Absyn.IdentDict

    fun procLabel id = "P" ^ Absyn.identName id
    fun varLabel id = "V" ^ Absyn.identName id

    exception AbsynToRTL
    fun bug msg =
      (TextIO.output(TextIO.stdErr, "Compiler error: "^msg^"\n");
       raise AbsynToRTL)

    (* XXX: REPLACE WITH YOUR CODE *)

    fun program(Absyn.PROGRAM{decs,...}) =
      RTL.PROGRAM([])

  end (* functor AbsynToRTLFn *)
