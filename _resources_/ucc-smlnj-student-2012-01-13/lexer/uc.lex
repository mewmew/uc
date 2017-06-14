(* lexer/uc.lex *)

(* parameters for ML-Yacc -- don't change *)
type arg = LexArg.lexarg
type pos = LexArg.pos
type svalue = Tokens.svalue
type ('a,'b) token = ('a,'b) Tokens.token
type lexresult = (svalue,pos) token

fun eof(lexarg) =
  let val pos = LexArg.readPos lexarg
  in
    (* complete? maybe or maybe not *)
    Tokens.EOF(pos, pos)
  end

(* YOUR HELPER FUNCTIONS HERE *)

(* YOUR HELPER DECLARATIONS BEFORE SECOND "double-percent" *)
(* YOUR TOKENS SPECIFICATION AFTER SECOND "double-percent" *)
(* sorry, but ml-lex doesn't allow comments in the sections below *)

%%

%header (functor UCLexFn(structure Tokens : UC_TOKENS
			 structure LexArg : LEXARG) : ARG_LEXER);
%arg (lexarg);
%full

%%

<INITIAL>"!="		=>
	(Tokens.NOTEQ(yypos, yypos+1));
