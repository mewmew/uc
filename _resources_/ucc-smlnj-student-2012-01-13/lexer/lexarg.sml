(* lexer/lexarg.sml *)

functor LexArgFn(Source : SOURCE) : LEXARG =
  struct

    structure Source = Source

    type pos = int
    datatype lexarg
      = A of {	fileName	: string,
		newLines	: int list ref,
		thisLine	: pos ref,
		leftPos		: pos ref,
		errFlag		: bool ref,
		readPos		: int ref }

    fun new(fileName, is) =
      let val readPos = ref 2	(*XXX: ML-Lex*)
	  fun yyinput n =
	    let val string = TextIO.inputN(is, n)
		val _ = readPos := !readPos + String.size string
	    in
	      string
	    end
	  val lexarg =
	    A{fileName = fileName,
	      newLines = ref [],
	      thisLine = ref 2,
	      leftPos = ref 0,
	      errFlag = ref false,
	      readPos = readPos}
      in
	(lexarg,yyinput)
      end

    fun newLine(A{newLines,thisLine,...}, pos) =
      (newLines := pos :: !newLines; thisLine := pos+1)

    fun newTab(A{thisLine,readPos,...}, yygone, yypos) =
      let val lpos = yypos - !thisLine
	  val incr = 7 - Int.rem(lpos, 8)
      in
	readPos := !readPos + incr;
	yygone := !yygone + incr
      end

    fun leftPos(A{leftPos,...}) = leftPos
    fun readPos(A{readPos,...}) = !readPos
    fun seenErr(A{errFlag,...}) = !errFlag

    fun source(A{fileName,newLines,...}) =
      Source.SOURCE{fileName = fileName, newLines = !newLines}

    fun error2 (lexarg as A{errFlag,...}) (msg,left,right) =
      (errFlag := true;
       Source.sayMsg (source lexarg) ("Error: "^msg, left, right))

  end (* functor LexArgFn *)
