(* main/main.sml *)

signature MAIN =
  sig
    val main: string list -> unit
  end (* signature MAIN *)

functor MainFn(structure Parse : PARSE
	       structure AbsynCheck : ABSYN_CHECK
	       structure AbsynPrint : ABSYN_PRINT
	       structure AbsynToRTL : ABSYN_TO_RTL
	       structure RTLPrint : RTL_PRINT
	       structure Codegen : CODEGEN
	       structure AssemPrint : ASSEM_PRINT
	       sharing Parse.Absyn = AbsynCheck.Absyn
	       sharing Parse.Absyn = AbsynPrint.Absyn
	       sharing Parse.Absyn = AbsynToRTL.Absyn
	       sharing RTLPrint.RTL = AbsynToRTL.RTL
	       sharing Codegen.RTL = AbsynToRTL.RTL
	       sharing Codegen.Assem = AssemPrint.Assem
		 ) : MAIN =
  struct

    fun sayErr msg = TextIO.output(TextIO.stdErr, msg)

    datatype 'a outcome = OK of 'a | ERR of exn

    fun withOutput f arg2 file =
      let val os = TextIO.openOut file
	  val outcome = (OK(f(os, arg2))) handle exn => ERR exn
      in
	TextIO.closeOut os;
	case outcome of (OK result) => result | (ERR exn) => raise exn
      end

    fun parse prefix = Parse.program(prefix ^ ".c")

    val printAbsyn = ref false
    val printRTL = ref false

    fun translate base =
      let val {dir=_,file} = OS.Path.splitDirFile base
	  val absyn = parse base
	  val _ = AbsynCheck.program absyn
	  val _ =
	    if !printAbsyn then
	      withOutput AbsynPrint.program absyn (file ^ ".absyn")
	    else ()
	  val rtl = AbsynToRTL.program absyn
	  val _ =
	    if !printRTL then
	      withOutput RTLPrint.program rtl (file ^ ".rtl")
	    else ()
	  val assem = Codegen.program rtl
	  val _ =
	    withOutput AssemPrint.program assem (file ^ AssemPrint.fileSuffix)
      in
	()
      end

    exception Usage
    fun usage badarg =
      (sayErr("ucc: invalid argument '" ^ badarg ^ "'\n");
       sayErr "usage: ucc [options] <file>.c ...\n";
       sayErr "available options:\n";
       sayErr "--print-absyn\n";
       sayErr "--print-rtl\n";
       raise Usage)

    fun option arg =
      case arg
	of "--print-absyn" => printAbsyn := true
	 | "--print-rtl" => printRTL := true
	 | _ => usage arg

    fun main argv =
      let fun process arg =
	    if String.sub(arg, 0) = #"-" then option arg
	    else
	      let val {base,ext} = OS.Path.splitBaseExt arg
	      in
		case ext
		  of SOME "c" => translate base
		   | _ => usage arg
	      end
      in
	List.app process argv
      end

  end (* functor MainFn *)
