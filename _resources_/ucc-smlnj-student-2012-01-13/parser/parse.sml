(* parser/parse.sml *)

functor ParseFn(structure Absyn : ABSYN
		structure Parser : ARG_PARSER
		  where type arg = unit
			and type pos = int
		structure LexArg : LEXARG
		sharing type Parser.result = Absyn.program
		sharing type Parser.lexarg = LexArg.lexarg
		sharing Absyn.Source = LexArg.Source
		  ) : PARSE =
  struct
    structure Absyn = Absyn

    fun program file =
      let val is = TextIO.openIn file
      in
	(let val (lexarg,inputf) = LexArg.new(file, is)
	     val lexer = Parser.makeLexer inputf lexarg
	     val (result,_) = Parser.parse(0,lexer,LexArg.error2 lexarg,())
	 in
	   if LexArg.seenErr lexarg then raise Parser.ParseError
	   else
	     let val _ = TextIO.closeIn is
		 val Absyn.PROGRAM{decs,...} = result
	     in
	       Absyn.PROGRAM{decs=decs,source=LexArg.source lexarg}
	     end
	 end) handle e => (TextIO.closeIn is; raise e)
      end

  end (* functor ParseFn *)
