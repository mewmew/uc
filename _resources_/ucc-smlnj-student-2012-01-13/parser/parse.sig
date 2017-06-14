(* parser/parse.sig *)

signature PARSE =
  sig

    structure Absyn	: ABSYN
    val program		: string -> Absyn.program

  end (* signature PARSE *)
