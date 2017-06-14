(* lexer/lexarg.sig *)

signature LEXARG =
  sig

    structure Source: SOURCE

    type pos = int
    type lexarg

    val new	: string * TextIO.instream -> lexarg * (int -> string)
    val newLine	: lexarg * pos -> unit
    val newTab	: lexarg * pos ref * pos -> unit
    val leftPos	: lexarg -> pos ref
    val readPos	: lexarg -> pos
    val seenErr	: lexarg -> bool
    val error2	: lexarg -> string * pos * pos -> unit
    val source	: lexarg -> Source.source

  end (* signature LEXARG *)
