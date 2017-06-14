(* parser/fake-parse.sml *)
(* This is only to be used for Assignment 1. *)

signature UC_TOKENS =
  sig
    type svalue (*= unit*)
    datatype ('a,'b) token
      = EOF of 'b * 'b
      | NOTEQ of 'b * 'b
      (* YOU NEED TO ADD REMAINING TOKENS HERE FOR ASSIGNMENT 1 *)

    val printToken: (svalue,int) token -> unit

  end (* signature UC_TOKENS *)

structure FakeTokens : UC_TOKENS =
  struct
    type svalue = unit
    datatype ('a,'b) token
      = EOF of 'b * 'b
      | NOTEQ of 'b * 'b
      (* YOU NEED TO ADD REMAINING TOKENS HERE FOR ASSIGNMENT 1 *)

    fun printToken(t) = ()	(* XXX: YOUR CODE *)
  end (* structure Tokens *)

functor FakeParseFn(structure Absyn : ABSYN
		    structure Lex : ARG_LEXER
		      where type UserDeclarations.pos = int
		    structure LexArg : LEXARG
		    structure FakeTokens : UC_TOKENS
		    sharing type Lex.UserDeclarations.token = FakeTokens.token
		    sharing type Lex.UserDeclarations.svalue = FakeTokens.svalue
		    sharing type LexArg.lexarg = Lex.UserDeclarations.arg
		    sharing Absyn.Source = LexArg.Source
		      ) : PARSE =
  struct

    structure Absyn = Absyn

    exception ParseError

    fun processTokens(lexer) = ()	(* XXX: YOUR CODE *)

    fun program file =
      let val is = TextIO.openIn file
      in
	(let val (lexarg,inputf) = LexArg.new(file, is)
	     val lexer = Lex.makeLexer inputf lexarg
	     val _ = processTokens(lexer)
	 in
	   if LexArg.seenErr lexarg then raise ParseError
	   else
	     let val _ = TextIO.closeIn is
	     in
	       Absyn.PROGRAM{decs=[], source=LexArg.Source.dummy}
	     end
	 end) handle e => (TextIO.closeIn is; raise e)
      end

  end (* functor FakeParseFn *)
