(* mips/mips.sml *)

(*
   This file implements the mips instruction set, 
   and SPIM assembler.

   For further information on see:
     http://www.cs.wisc.edu/~larus/SPIM/cod-appa.pdf
     Assemblers, Linkers, and the SPIM Simulator
      in Hennessy & Patterson, 
	 Computer Organization and Design: 
         The Hardware/Software Interface.  
   References to pages below are to this book.

   Please report any bugs to happi@csd.uu.se.

   Modifications 2003-09-08 by Mikael Pettersson (mikpe@csd.uu.se):
   1. Replaced Symbol.symbol by string.
   2. Replaced ErrorMsg.impossible by local error function.
   3. Split up code union to separate move/oper/jump types,
      and added new insn union. This makes MIPS compatible
      with the INSNS signature, except the *DefUse and jumpTargets
      functions aren't yet implemented.
   4. Added isPseudoReg/isPseudoFReg predicates.

   Modifications 2003-12-01 by Mikael Pettersson (mikpe@csd.uu.se):
   1. Corrected byte and halfword load instructions to allow
      general effective addresses, not just labels.
   2. Added missing byte and halfword store instructions.
*)



structure MIPS : MIPS =
  struct

    type reg = int	(* integer register (5-bits 0-31)      *)
    type freg = int	(* floatingpoint register (5-bits 0-31)*)

	                (* Some predefined registers: (p. A-23)*)
    val RV  = 02        (* Return value                        *)
    val FRV = 00        (* Float return value                  *)
    val SP  = 29        (* Stack pointer                       *)
    val FP  = 30        (* Frame pointer                       *)
    val RA  = 31        (* Return address                      *)


    type imm = int      (* Immediate integer                   *)
    type neg = bool     (* Flag for negative offsets           *)
    type label = string (* A label...                  *)  
	
    datatype ea 	(* Effective address (pp. A-49 - A-50) *)
      = IMM of imm                    (* Immediate             *)
      | REG of reg                    (* Contents of register  *)
      | OFFSET of reg * imm           (* Imm + cont. of reg    *)
      | LAB of label                  (* Address of label      *)
      | LABOFF of label * neg * imm   (* -||- +/- imm          *)
      | LABINDOFF of label * neg * 
                     imm * reg        (* -||- +/- (imm + reg)  *)

    datatype aluop =    (* Arithmetic logic ops r1 = r2 op r3  *)
	(* Arithmetic *)
	ADD             (* add (with overflow)       (p. A-55) *)
      | ADDU            (* add (without overflow)    (p. A-55) *)
      | AND             (* and                       (p. A-55) *) 
      | DIV             (* div (with overflow)       (p. A-56) *) 
      | DIVU            (* div (without overflow)    (p. A-56) *) 
      | MUL             (* multiply(without overflow)(p. A-56) *) 
      | MULO            (* multiply (with overflow)  (p. A-56) *) 
      | MULOU           (* unsigned multiply 
                                   (without overflow)(p. A-57) *) 
      | REM             (* reminider                 (p. A-57) *)  
      | REMU            (* unsigned reminider        (p. A-58) *)  
      | SUB             (* subtract (with owerflow)  (p. A-59) *)
      | SUBU            (* subtract(without owerflow)(p. A-59) *)

	(* Logic       *)
      | NOR             (* nor                       (p. A-57) *)  
      | OR              (* or                        (p. A-57) *)  
      | XOR             (* xor                       (p. A-59) *)  
	
      | SLLV            (* shift left logical        (p. A-58) *)  
      | SRAV            (* shift right arithmetic    (p. A-58) *)
      | SRLV            (* shift right logical       (p. A-58) *) 
      | ROL             (* rotate left               (p. A-58) *)
      | ROR             (* rotate right              (p. A-59) *)

	(* Comparison   if r2 OP r3 then r1 = 1 else r1 = 0    *)
      | SLT             (* set less than             (p. A-60) *)
      | SLTU            (* set less than unsigned    (p. A-60) *)
      | SEQ             (* set equal                 (p. A-60) *)
      | SGE             (* set greater than or equal (p. A-60) *)
      | SGEU            (* -||- unsigned             (p. A-60) *)
      | SGT             (* set greater than          (p. A-61) *)
      | SGTU            (* -||- unsigned             (p. A-61) *)
      | SLE             (* set less than or equal    (p. A-61) *)
      | SLEU            (* -||- unsigned             (p. A-61) *)
      | SNE             (* set not equal             (p. A-61) *)

    datatype immop =    (* r1 = r2 OP imm                      *)
	ADDI            (* add (with overflow)       (p. A-55) *)
      | ADDIU           (* add (without overflow)    (p. A-55) *)

      | ANDI            (* and                       (p. A-55) *)
      | XORI            (* xor                       (p. A-59) *) 
      | ORI             (* or                        (p. A-57) *)

      | SLTI            (* set less than             (p. A-60) *)
      | SLTIU           (* -||- unsigned             (p. A-60) *)

    datatype binop =    (* r1 = OP r2                          *)
	ABS             (* absolute value            (p. A-55) *)
      | NEG             (* negate value              (p. A-57) *)
      | NEGU            (* -||- without overflow     (p. A-57) *)
      | NOT             (* not                       (p. A-57) *)

    datatype isop =     (* shift immediate r1 = r2 SOP imm     *)
	SLL             (* shift left logical        (p. A-58) *)  
      | SRA             (* shift right arithmetic    (p. A-58) *) 
      | SRL             (* shift right logical       (p. A-58) *) 

    datatype fop =      (* Floating point op fr1 = fr2 FOP fr3 *)
	                (*  single precision                   *)
	FADD            (* add                       (p. A-70) *)  
      | FSUB            (* sub                       (p. A-74) *) 
      | FMUL            (* multiply                  (p. A-73) *) 
      | FDIV            (* divide                    (p. A-72) *) 

    datatype fcomp =    (* Floating point comparison           *)
	                (* set fp-flag if fr1 COP fr2          *)
	                (* Use bc1t or bc1f to test the result *)
	FEQ             (* fr1 = fr2                 (p. A-71) *)
      | FLE             (* fr1 =< fr2                (p. A-71) *)
      | FLT             (* fr1 < fr2                 (p. A-71) *)	

    datatype ubranch =  (* Branch on OP r1                     *)
	BGEZ            (* r1 >= 0 -> goto LAB       (p. A-62) *)
      | BGTZ            (* r1 > 0 -> goto LAB        (p. A-62) *)
      | BLEZ            (* r1 =< 0 -> goto LAB       (p. A-63) *)
      | BLTZ            (* r1 < 0 -> goto LAB        (p. A-63) *)
      | BEQZ            (* r1 == 0 -> goto LAB       (p. A-63) *)
      |	BNEZ            (* r1 <> 0 -> goto LAB       (p. A-64) *)

	                (* The following instructions are really call
                           instructions, the return address register
                           is set to the next instruction.     *)
      | BEGEZAL         (* r1 >= 0 -> call LAB       (p. A-62) *)
      | BLTZAL          (* r1 < 0 -> call LAB        (p. A-63) *)

    datatype bbranch =  (* Branch if r1 COMP r2                *)
	BEQ             (* =  *)
      | BNE             (* <> *)
      | BGE             (* >= *)
      | BGT             (* >  *)
      | BLE             (* <= *)
      | BLT             (* <  *)
      | BGEU            (* >= unsigned *)
      | BGTU            (* > unsigned  *)
      | BLEU            (* <= unsigned *)
      |	BLTU            (* < unsigned  *)

    (* Don't use this constructors directly if not absolutely neccessary,
       use the primitives below instead. *)

    datatype move			 (* reg-to-reg copies	   *)
      = IMOVE of reg * reg                (* r1 := r2               *) 
      | FMOVE of freg * freg             (* fr1 := fr2             *) 

    datatype oper
      = BOP of binop * reg * reg         (* ALU op r1 = OP r2      *)
      | OP of aluop * reg * reg * reg    (* ALU op r1 = r2 OP r3   *)
      | IOP of immop * reg * reg * imm   (* ALU op r1 = r2 OP imm  *)
      | ISOP of isop * reg * reg * imm   (* Shift  r1 = r2 OP imm  *)
      | FOP of fop * freg * freg *freg   (* FP op fr1 = fr2 OP fr3 *)
      | FCOMP of fcomp * freg * freg     (* Comp Fflag = fr1 OP fr2 *)
      | LI of reg * imm                  (* r1 := imm              *)
      | LA of reg * label                (* Load address (p. A-65) *)
      | LEA of reg * ea                  (* Load effective address *)
      | LB of reg * ea                   (* Load byte              *)
      | LBU of reg * ea                  (* Load unsigned byte     *)
      | LH of reg * ea                   (* Load halfword          *)
      | LHU of reg * ea                  (* Load unsigned halfword *)
      | LW of reg * ea                   (* Load word              *)
      | SB of reg * ea                   (* Store byte             *)
      | SH of reg * ea                   (* Store halfword         *)
      | SW of reg * ea                   (* Store word             *)
      | L_D of freg * ea                 (* Load double prec.float *)
      | L_S of freg * ea                 (* Load float             *)
      | S_S of freg * ea                 (* Store float            *)
      | S_D of freg * ea                 (* Store double           *)
      | MFC1 of reg * freg        (* Move from coprocessor 1       *)
      | MTC1 of freg * reg        (* Move to coprocessor 1         *)
      | CVTSW of freg * freg      (* Convert int to float          *)
      | CVTWD of freg * freg      (* Convert double to int         *)
      | CVTWS of freg * freg      (* Convert float to int          *)
      | SYSCALL                   (* System call (pp. A-48 - A-49) *)  
      | BREAK of imm              (* Casue exception code.         *)
      | NOP                       (* No operation                  *)

    datatype jump
      = B of label                       (* Goto label             *)
      | B1 of ubranch * reg * label      (* If OP(r1) -> B label   *)   
      | B2 of bbranch * reg * reg * label(* If r1 OP r2 -> B label *)
	                                 (* Coprocessor 1 is FPU   *)
	                                 (* Set F-flag with FCOMP
					    instruction, then branch
					    with these isntructions *)
      | BCZT of int * label              (* If coprocessor Z-flag 
					    -> label               *)
      | BCZF of int * label              (* If not coprocessor Z-flag 
					    -> label               *)
      | JAL of ea                 (* JumpAndLink (Call)  (p. A-65) *)
      | J of ea                   (* Goto effective address        *)
      | JALR of reg * reg         (* Like JAL but store ret.address in r2 *)
      | RETURN                    (* Jump to address in RA         *)

    datatype insn
      = LABEL of label                   (* Name a point in code   *)  
      | MOVE of move
      | OPER of oper
      | JUMP of jump

    exception MIPS
    fun error(msg) =
      (TextIO.output(TextIO.stdErr, "Compiler error: "^msg^"\n");
       raise MIPS)

   (* Interface functions to the mips instructions *)
  
   (* Immediate aluops, reg1 := reg2 op imm (where imm is an integer) *)
    fun subi(reg1, reg2, imm) = OPER(IOP(ADDI, reg1, reg2, ~imm))
    fun addi(reg1, reg2, imm) = OPER(IOP(ADDI, reg1, reg2, imm))

    fun fadd(reg1, reg2, reg3) = OPER(FOP(FADD, reg1, reg2, reg3))

    fun load_imm(reg, imm) = OPER(LI(reg,imm))
    fun load_float(freg, address) = OPER(L_S(freg,LAB(address)))
    fun store_float(freg, address) = OPER(S_S(freg,LAB(address)))
    fun store(reg1, reg2, offset) = OPER(SW(reg1,OFFSET(reg2,offset)))
    fun load(reg1, reg2, offset) = OPER(LW(reg1,OFFSET(reg2,offset)))
	
    fun b_on_fflag_true lab = JUMP(BCZT(1,lab))
    fun b_on_fflag_false lab = JUMP(BCZF(1,lab))
    fun label l = LABEL l
    fun call(address) = JUMP(JAL(LAB(address)))
    fun jump(address) = JUMP(J(LAB(address)))
    fun jump_reg(r) = JUMP(J(REG(r)))
    val return = JUMP(RETURN)
    val syscall = OPER(SYSCALL)

    datatype asm 
	= COMMENT of string           (* A one line comment *)
      | INSTRUCTION of insn           (* A machine instruction *)    
      | COMINST of insn * string      (* A commented instruction *)

      | DATA of label                 (* Named data *)

      | ALLOC of label * int          (* Name data of size *)
      | DOUBLE of label * real        (* Named double *)

      | ALIGN of int                  (* Align next datum to 2^n *)
      | ASCII of string               (* Store the ascii string  *)
      | ASCIIZ of string              (* 0-terminated string  *)
      | BYTE of int list              (* Store bytes *)
      | DATASEG                       (* Start data segment *)
      | DATASEG_AT of int             (* Put data at address *)    
      | DOUBLES of real list          (* Store floats in mem *)
      | EXTERN of label * int         (* Declare external data *)
      | FLOAT of label * real         (* A named float *)
      | FLOATS of real list           (* Single precision numbers *)
      | GLOBAL of label               (* Declare global data *)  
      | HALF of int list              (* Store 16-bit halfwords *)
      | KDATA                         (* Start Kernel segment *)
      | KDATA_AT of int               (*  - " - at address *) 
      | SPACE of int                  (* Allocate n bytes *)
      | TEXT                          (* Start Text segment *)
      | TEXT_AT of int                (*  - " - at address *) 
      | WORD of int list              (* Store 32-bit words *)
      
   (* Make room for a var or array, and give it a name. *)
    fun alloc (name, size) = ALLOC(name, size)
    fun alloc_float (name,f) = FLOAT(name, f)
    fun align n = ALIGN(n)

    type program = asm list          (* A list of instructions *)


    fun print_program(program:program, outstream) =
	let

	    fun say s = TextIO.output (outstream,s)
	    fun fills os (pos, limit) =
		if pos < limit then 
		    (TextIO.output (os," ");
		     fills os (pos+1, limit))
		else
		    pos
		    
	    val LPOS = 0
	    val OPOS = 10
	    val APOS = 20
	    val CPOS = 40
	    
	    val fill = fills outstream
	    fun toLower s = 
		implode (List.map (fn x => Char.toLower x) (explode s));
	    fun length s = String.size s
	    fun print s = (TextIO.output (outstream,s);length s)
	    fun print_n 0 = 0
	      | print_n 1 = (print "0")
	      | print_n n = 
		(print "0,") +
		print_n (n-1)
	    fun print_int(i) =
		let val s = 
		    if i >= 0 then Int.toString(i)
		    else "-" ^ Int.toString(~i)
		in
		    print s
		end
	    fun print_float(f) =
		let val s = 
		    if f >= 0.0 then Real.toString(f)
		    else "-" ^ Real.toString(~f)
		in
		    print s
		end
	    fun print_byte(byte) =
		if byte <= 255 andalso byte >= 0 then
		    print_int(byte)
		else
		    error("too large byte: "^Int.toString byte)
		    
	    fun print_list f [e] = 
		f e
	      | print_list f (e::es) =
		f e + 
		print ", " +
		print_list f es
	      | print_list f [] = 0

	    fun print_bytes(bs) = 
		print_list print_byte bs
		
	    fun print_hword(hword) =
		if hword <= 65535 andalso hword >= 0 then
		    print_int(hword)
		else
		    error("too large hword: "^Int.toString hword)
		    
	    fun print_hwords(hws) = 
		print_list print_hword hws

	    fun print_words(ws) = 
		print_list print_int ws

	    fun print_floats(fs) =
		print_list print_float fs
		
	    fun pi(instr) =
		fill(fill (LPOS,OPOS) + print(instr),APOS) 
	      
	    fun print_freg(r) =
		print "$f" + print_int(r)
		
	    fun print_reg(r) =
		print (case r of
		       0 => "$zero"  (* Constant 0             *)
		     | 1 => "$at"    (* Reserved for assembler *) 
		     | 2 => "$v0"    (* Return value 1         *)
		     | 3 => "$v1"    (* Return value 2         *)
		     | 4 => "$a0"    (* Argument 1             *)
		     | 5 => "$a1"    (* Argument 2             *)
		     | 6 => "$a2"    (* Argument 3             *)
		     | 7 => "$a3"    (* Argument 4             *)
		     | 8 => "$t0"    (* Temp  1  - caller-save *)  
		     | 9 => "$t1"    (* Temp  2      --||--    *)  
		     |10 => "$t2"    (* Temp  3      --||--    *) 
		     |11 => "$t3"    (* Temp  4      --||--    *) 
		     |12 => "$t4"    (* Temp  5      --||--    *)  
		     |13 => "$t5"    (* Temp  6      --||--    *)  
		     |14 => "$t6"    (* Temp  7      --||--    *)  		     
		     |15 => "$t7"    (* Temp  8      --||--    *) 
		     |16 => "$s0"    (* Temp  1s - callee-save *)  
		     |17 => "$s1"    (* Temp  2s     --||--    *)  
		     |18 => "$s2"    (* Temp  3s     --||--    *)
		     |19 => "$s3"    (* Temp  4s     --||--    *)
		     |20 => "$s4"    (* Temp  5s     --||--    *)
		     |21 => "$s5"    (* Temp  6s     --||--    *)
		     |22 => "$s6"    (* Temp  7s     --||--    *)
		     |23 => "$s7"    (* Temp  8s     --||--    *)
		     |24 => "$t8"    (* Temp  9  - caller-save *)
		     |25 => "$t9"    (* Temp 10      --||--    *)
		     |26 => "$k0"    (* Reserved for OS        *) 
		     |27 => "$k1"    (* Reserved for OS        *) 
		     |28 => "$gp"    (* Pointer to global area *)
		     |29 => "$sp"    (* Stack pointer          *)
		     |30 => "$fp"    (* Frame pointer          *)
		     |31 => "$ra"    (* Return address	       *)
		     |_ => error("invalid register: "^Int.toString r))

	fun print_label lab = print (lab)
	fun print_ea_reg(r) = 
	    (print "(") + (print_reg r) + (print ")")
	fun print_ea(IMM i) = print_int(i)
	  | print_ea(REG r) = print_ea_reg r
	  | print_ea(OFFSET (r,i)) =
	    (print_int i) + print_ea_reg r
	  | print_ea(LAB s) = print_label s
	  | print_ea(LABOFF(lab,false,i)) =
	             (print_label lab) +
		     (print " + ") +  
		     (print_int i)
	  | print_ea (LABOFF(lab,true,i)) =
	             (print_label lab) +
		     (print " - ") +  
		     (print_int i)
	  | print_ea (LABINDOFF(lab,false,i,r)) =
	             (print_label lab) +
		     (print " + ") +  
		     (print_int i) +
		     (print_ea_reg r)
	  | print_ea (LABINDOFF(lab,true,i,r)) =
	             (print_label lab) +
		     (print " - ") + 
		     (print_int i) +
		     (print_ea_reg r)

	fun print_aluop(ADD)   = print("add")
	  | print_aluop(ADDU)  = print("addu")
	  | print_aluop(AND)   = print("and")
	  | print_aluop(DIV)   = print("div")
	  | print_aluop(DIVU)  = print("divu")
	  | print_aluop(MUL)   = print("mul")
	  | print_aluop(MULO)  = print("mulo")
	  | print_aluop(MULOU) = print("mulou")
	  | print_aluop(NOR)   = print("nor")
	  | print_aluop(OR)    = print("or")
	  | print_aluop(REM)   = print("rem")
	  | print_aluop(REMU)  = print("remu")
	  | print_aluop(SLLV)  = print("sllv")
	  | print_aluop(SRAV)  = print("srav")
	  | print_aluop(SRLV)  = print("srlv")
	  | print_aluop(ROL)   = print("rol")
	  | print_aluop(ROR)   = print("ror")
	  | print_aluop(SUB)   = print("sub")
	  | print_aluop(SUBU)  = print("subu")
	  | print_aluop(XOR)   = print("xor")
	  | print_aluop(SLT)   = print("slt")
	  | print_aluop(SLTU)  = print("sltu")
	  | print_aluop(SEQ)   = print("seq")
	  | print_aluop(SGE)   = print("sge")
	  | print_aluop(SGEU)  = print("sgeu")
	  | print_aluop(SGT)   = print("sgt")
	  | print_aluop(SGTU)  = print("sgtu")
	  | print_aluop(SLE)   = print("sle")
	  | print_aluop(SLEU)  = print("sleu")
	  | print_aluop(SNE)   = print("sne")

	fun print_immop(ADDI)  = print("addi")
	  | print_immop(ADDIU) = print("addiu")
	  | print_immop(ANDI)  = print("andi")
	  | print_immop(ORI)   = print("ori")
	  | print_immop(XORI)  = print("xori")
	  | print_immop(SLTI)  = print("slti")
	  | print_immop(SLTIU) = print("sltiu")
	    
	fun print_binop(ABS)  = print("abs")
	  | print_binop(NEG)  = print("neg")
	  | print_binop(NEGU) = print("negu")
	  | print_binop(NOT)  = print("not")

	fun print_isop(SLL)   = print("sll")
	  | print_isop(SRA)   = print("sra")
	  | print_isop(SRL)   = print("srl")

	    


	    
        fun print_ub(BGEZ)    = print("bgez")
	  | print_ub(BEGEZAL) = print("begezal")
	  | print_ub(BGTZ)    = print("bgtz")
	  | print_ub(BLEZ)    = print("blez")
	  | print_ub(BLTZAL)  = print("bltzal")
	  | print_ub(BLTZ)    = print("bltz")
	  | print_ub(BEQZ)    = print("beqz")
	  | print_ub(BNEZ)    = print("bnez")

	fun print_bb(BEQ)    = print("beq")
	  | print_bb(BNE)    = print("bne")
	  | print_bb(BGE)    = print("bge")
	  | print_bb(BGEU)   = print("bgeu")
	  | print_bb(BGT)    = print("bgt")
	  | print_bb(BGTU)   = print("bgtu")
	  | print_bb(BLE)    = print("ble")
	  | print_bb(BLEU)   = print("bleu")
	  | print_bb(BLT)    = print("blt")
	  | print_bb(BLTU)   = print("bltu")


	fun print_fop(FADD) = print("add.s")
	  | print_fop(FSUB) = print("sub.s")
	  | print_fop(FMUL) = print("mul.s")
	  | print_fop(FDIV) = print("div.s")

	fun print_fcomp(FEQ) = print("c.eq.s")
	  | print_fcomp(FLE) = print("c.le.s")
	  | print_fcomp(FLT) = print("c.lt.s")

	fun printLabel(lab) = (print_label lab) + print(":")

	fun printMove(move) =
	  case move
	    of IMOVE(r1,r2) =>
		    fill(fill(0,OPOS) + 
			 print "move",
			 APOS) +
		    (print_reg r1) +
		    print(", ") + (print_reg r2)
	      | FMOVE(r1,r2) =>
		    fill(fill(0,OPOS) + 
			 print "mov.s",
			 APOS) +
		    (print_freg r1) +
		    print(", ") + (print_freg r2)

	fun printOper(oper) =
	    case oper of
		BOP(bop,r1,r2) =>
		    fill(fill(0,OPOS) + 
			 print_binop(bop),APOS) +
		    (print_reg r1) +
		    print(", ") + (print_reg r2) 
	      | OP(aluop, r1, r2, r3) =>
		    fill(fill(0,OPOS) + 
			 print_aluop(aluop),APOS) +
		    (print_reg r1) +
		    print(", ") + (print_reg r2) +
		    print(", ") + (print_reg r3) 
	      | IOP(immop, r1, r2, i) =>
		    fill(fill(0,OPOS) + 
			 print_immop(immop),
			 APOS) +
		    (print_reg r1) +
		    print(", ") + (print_reg r2) +
		    print(", ") + (print_int i) 
	      | ISOP(isop, r1, r2, i) =>
		    fill(fill(0,OPOS) + 
			 print_isop(isop),
			 APOS) +
		    (print_reg r1) +
		    print(", ") + (print_reg r2) +
		    print(", ") + (print_int i) 
	      | FOP(fop,r1,r2,r3) =>
		    fill(fill(0,OPOS) + 
		    print_fop(fop),APOS) +
		    (print_freg r1) +
		    print(", ") + (print_freg r2) +
		    print(", ") + (print_freg r3) 
	      | FCOMP(fcomp,r1,r2) =>
		    fill(fill(0,OPOS) + 
		    print_fcomp(fcomp),APOS) +
		    (print_freg r1) +
		    print(", ") + (print_freg r2)
	      | LA(r,l) => 
		    fill(pi("la"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_label l)
	      | LEA(r,address) => 
		    fill(pi("la"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_ea address)
	      | LB(r,address) => 
		    fill(pi("lb"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_ea address)
	      | LBU(r,address) =>
		    fill(pi("lbu"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_ea address)
	      | LH(r,address) =>
		    fill(pi("lh"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_ea address)
	      | LHU(r,address) =>
		    fill(pi("lhu"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_ea address)
	      | LW(r,address) =>
		    fill(pi("lw"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_ea address)
	      | SB(r,address) =>
		    fill(pi("sb"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_ea address)
	      | SH(r,address) =>
		    fill(pi("sh"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_ea address)
	      | SW(r,address) =>
		    fill(pi("sw"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_ea address)
	      | L_D(fr,address) =>
		    fill(pi("l.d"),APOS) +
		    print_freg(fr) +
		    print(", ") + (print_ea address)
	      | L_S(fr,address) =>
		    fill(pi("l.s"),APOS) +
		    print_freg(fr) +
		    print(", ") + (print_ea address)
	      | LI(r,i) =>
		    fill(pi("li"),APOS) +
		    print_reg(r) +
		    print(", ") + (print_int i)
	      | S_S(fr,address) =>
		    fill(pi("s.s"),APOS) +
		    print_freg(fr) +
		    print(", ") + (print_ea address)
	      | S_D(fr,address) =>
		    fill(pi("s.d"),APOS) +
		    print_freg(fr) +
		    print(", ") + (print_ea address)
	      | MFC1(r,fr) =>
		    fill(pi("mfc1"),APOS) +
		    print_reg(r) +
		    print(", ") + 
		    print_freg(fr)
	      | MTC1(fr,r) =>
		    fill(pi("mtc1"),APOS) +
		    print_reg(r) +
		    print(", ") + 
		    print_freg(fr)
	      | CVTSW(fr1,fr2) =>
		    fill(pi("cvt.s.w"),APOS) +
		    print_freg(fr1) +
		    print(", ") + 
		    print_freg(fr2)
	      | CVTWD(fr1,fr2) =>
		    fill(pi("cvt.w.d"),APOS) +
		    print_freg(fr1) +
		    print(", ") + 
		    print_freg(fr2)
	      | CVTWS(fr1,fr2) =>
		    fill(pi("cvt.w.s"),APOS) +
		    print_freg(fr1) +
		    print(", ") + 
		    print_freg(fr2)
	      | SYSCALL => 
		    fill(pi("syscall"),APOS)
	      | BREAK(code) =>
		    fill(pi("break"),APOS) +
		    print_int(code)
	      | NOP =>
		    fill(pi("nop"),APOS)

	fun printJump(jump) =
	  case jump
	    of B(label) => 
		    fill(fill(0,OPOS) + 
			 print("b"),APOS) +
		    (print_label label)
	      | B1(ubranch, r1, label) =>
		    fill(fill(0,OPOS) + 
			 print_ub(ubranch),APOS) +
		    (print_reg r1) +
		    print(", ") + (print_label label)
	      | B2(bbranch, r1, r2, label) =>
		    fill(fill(0,OPOS) + 
			 print_bb(bbranch),APOS) +
		    (print_reg r1) +
		    print(", ") + (print_reg r2) +
		    print(", ") + (print_label label)
	      | BCZT(n, lab) =>
		    fill(fill(0,OPOS) + 
			 print("bc1t"),APOS) +
		    (print_label lab)
	      | BCZF(n, lab) =>
		    fill(fill(0,OPOS) + 
			 print("bc1f"),APOS) +
		    (print_label lab)
	      | JAL address =>
		    fill(pi("jal"),APOS) +
		    (print_ea address)
	      | J address =>
		    fill(pi("j"),APOS) +
		    (print_ea address)
	      | JALR(r1, r2) =>
		    fill(fill(0,OPOS) + 
			 print("jalr"),APOS) +
		    (print_reg r1) +
		    print(", ") + (print_reg r2)
	      | RETURN => 
		    fill(pi("jr"),APOS)+ print("$ra")

	fun printInsn(insn) =
	  case insn
	    of LABEL lab => printLabel(lab)
	     | MOVE move => printMove(move)
	     | OPER oper => printOper(oper)
	     | JUMP jump => printJump(jump)

	fun print_asm i = 
	    case i of
		COMMENT comm =>
		    say ("# " ^ comm ^ "\n")
	      | INSTRUCTION insn =>
		    (printInsn(insn);
		     say "\n")
	      | COMINST(insn, comm) =>
		    let 
			val pos = printInsn(insn)
		    in
			fill(pos,CPOS);
			say ("#" ^ comm ^ "\n")
		    end
	      | GLOBAL sym => say(".globl " ^ (sym) ^ "\n")
	      | DATA sym => say(".data " ^ (sym) ^ "\n")
	      | DOUBLE(id,f) =>
		    (fill(fill((print_label id) + (print ":"),OPOS) +
			 print(".double "), APOS)+
		    print_float(f);say("\n"))
	      | ALLOC(id,size) =>
		    (fill(fill((print_label id) + (print ":"),OPOS) +
			 print(".word "), APOS)+
		     print_n(size);say("\n"))
	      | ALIGN n => (say(".align "); print_int(n); say "\n")
	      | ASCII s => say (".ascii \"" ^ s ^ "\"\n")               
	      | ASCIIZ s => say (".asciiz \"" ^ s ^ "\"\n") 
	      | BYTE bytes => (say (".byte ");
			       print_bytes(bytes);
			       say "\n")
	      | DATASEG => say(".data\n")                      
	      | DATASEG_AT adr => (say(".data "); print_int(adr); say "\n")
	      | DOUBLES floats => (say(".double ");
				   print_floats(floats);
				   say "\n")
	      | EXTERN (lab,size) => (say ".extern ";
				      (print_label lab);
				      say " ";
				      print_int size;
				      say "\n")
	      | FLOAT (id,f) => 
		    (fill(fill((print_label id) + (print ":"),OPOS) +
			 print(".float "), APOS)+
		    print_float(f);say("\n"))

	      | FLOATS floats =>  (say(".float ");
				  print_floats(floats);
				  say "\n")       
	      | HALF hwords =>  (say (".half ");
				 print_hwords(hwords);
				 say "\n")             
	      | KDATA => say ".kdata\n"
	      | KDATA_AT adr => (say(".kdata "); print_int(adr); say "\n") 
	      | SPACE n =>   (say(".space "); print_int(n); say "\n")
	      | TEXT  => say ".text\n"                        
	      | TEXT_AT adr => (say(".text "); print_int(adr); say "\n")
 	      | WORD words =>  (say (".word ");
				 print_words(words);
				 say "\n")             
	in 
	    List.app print_asm program;
	    ()
	end

    fun isPseudoReg(reg) = reg > 31
    fun isPseudoFReg(freg) = freg > 31

  end (* structure MIPS *)
