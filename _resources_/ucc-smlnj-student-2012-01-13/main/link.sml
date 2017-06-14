(* main/link.sml *)

structure Absyn = AbsynFn(Source);

structure LexArg = LexArgFn(Source);

(* These two are for Assignment 1. Remove them afterwards. *)
structure UCLex =
  UCLexFn(structure Tokens = FakeTokens
	  structure LexArg = LexArg
	    );
structure FakeParse =
  FakeParseFn(structure Absyn = Absyn
	      structure Lex = UCLex
	      structure LexArg = LexArg
	      structure FakeTokens = FakeTokens
		);

(* These four must be commented out until you start Assignment 2. *)
(*
structure UCLrVals =
  UCLrValsFn(structure Token = LrParser.Token
	     structure Absyn = Absyn
	     structure LexArg = LexArg
	       );
structure UCLex =
  UCLexFn(structure Tokens = UCLrVals.Tokens
	  structure LexArg = LexArg
	    );
structure UCParser =
  JoinWithArg(structure Lex = UCLex
	      structure ParserData = UCLrVals.ParserData
	      structure LrParser = LrParser
		);
structure Parse =
  ParseFn(structure Absyn = Absyn
	  structure Parser = UCParser
	  structure LexArg = LexArg
	    );
*)

structure AbsynCheck = AbsynCheckFn(Absyn);
structure AbsynPrint = AbsynPrintFn(Absyn);
structure AbsynToRTL =
  AbsynToRTLFn(structure Absyn = Absyn
	       structure RTL = RTL
		 );
structure RTLPrint = RTLPrintFn(RTL);

structure MIPSCodegen =
  MIPSCodegenFn(structure RTL = RTL
		structure MIPS = MIPS
		  );
structure MIPSAssemPrint =
  MIPSAssemPrintFn(structure MIPS = MIPS
		     );

structure Main =
  MainFn(structure Parse = FakeParse (* XXX: change to Parse after Assignment 1 *)
	 structure AbsynCheck = AbsynCheck
	 structure AbsynPrint = AbsynPrint
	 structure AbsynToRTL = AbsynToRTL
	 structure RTLPrint = RTLPrint
	 structure Codegen = MIPSCodegen
	 structure AssemPrint = MIPSAssemPrint
	   );
