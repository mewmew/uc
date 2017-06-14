(* util/source.sig *)

signature SOURCE =
  sig

    datatype source
      = SOURCE of
	  { fileName: string,
	    newLines: int list }	(* _descending_ order *)

    val dummy	: source
    val sayMsg	: source -> string * int * int -> unit

  end (* signature SOURCE *)
