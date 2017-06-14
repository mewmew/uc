(* mips/mips.sig *)

signature MIPS =
  sig

    type reg = int
    type freg = int
    val RV : reg
    val FRV: reg
    val SP : reg
    val FP : reg
    val RA : reg

    type imm = int
    type neg = bool
    type label = string
	
    datatype ea 	
      = IMM of imm
      | REG of reg
      | OFFSET of reg * imm
      | LAB of label
      | LABOFF of label * neg * imm
      | LABINDOFF of label * neg * imm * reg


    datatype aluop = ADD | ADDU | AND | 
	DIV | DIVU | MUL | MULO | MULOU |
	NOR | OR | REM | REMU | SLLV | SRAV | 
	SRLV | ROL | ROR | SUB | SUBU |
	XOR | SLT | SLTU | SEQ | SGE | 
	SGEU | SGT | SGTU | SLE | SLEU | SNE
    datatype immop = ADDI | ADDIU | ANDI | ORI | XORI | SLTI | SLTIU
    datatype binop = ABS | NEG | NEGU | NOT
    datatype isop = SLL | SRA | SRL
    datatype fop = FADD | FSUB | FMUL | FDIV
    datatype fcomp = FEQ | FLE | FLT 	
    datatype ubranch = BGEZ | BEGEZAL | BGTZ | BLEZ | BLTZAL | BLTZ | BEQZ |
	BNEZ
    datatype bbranch = BEQ | BNE | BGE | BGEU | BGT | BGTU | BLE | BLEU | BLT |
	BLTU

    (* Don't use this constructors directly if not absolutely neccessary,
       use the primitives below instead.
       ... if they exist. *)
 
    datatype move
      = IMOVE of reg * reg
      | FMOVE of freg * freg

    datatype oper
      = BOP of binop * reg * reg
      | OP of aluop * reg * reg * reg
      | IOP of immop * reg * reg * imm
      | ISOP of isop * reg * reg * imm
      | FOP of fop * freg * freg *freg
      | FCOMP of fcomp * freg * freg
      | LI of reg * imm
      | LA of reg * label
      | LEA of reg * ea
      | LB of reg * ea
      | LBU of reg * ea
      | LH of reg * ea
      | LHU of reg * ea
      | LW of reg * ea
      | SB of reg * ea
      | SH of reg * ea
      | SW of reg * ea
      | L_D of freg * ea
      | L_S of freg * ea
      | S_S of freg * ea
      | S_D of freg * ea
      | MFC1 of reg * freg
      | MTC1 of freg * reg
      | CVTSW of freg * freg
      | CVTWD of freg * freg
      | CVTWS of freg * freg
      | SYSCALL
      | BREAK of imm
      | NOP 

    datatype jump
      = B of label
      | B1 of ubranch * reg * label
      | B2 of bbranch * reg * reg * label
      | BCZT of int * label
      | BCZF of int * label
      | JAL of ea
      | J of ea
      | JALR of reg * reg
      | RETURN

    datatype insn
      = LABEL of label
      | MOVE of move
      | OPER of oper
      | JUMP of jump

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
     


    type program = asm list          (* A list of instructions *)

    val label : label -> insn
    val addi : reg*reg*imm -> insn
    val subi : reg*reg*imm -> insn
    val fadd : freg*freg*freg -> insn
    val load_imm : reg*imm -> insn	
    val load_float : freg * label -> insn
    val store_float : freg * label -> insn
    val store : reg * reg * imm -> insn	
    val load : reg * reg * imm -> insn	
    val b_on_fflag_true : label -> insn
    val b_on_fflag_false : label -> insn

    val call : label -> insn 
    val jump : label -> insn 
    val jump_reg : reg -> insn 
    val return : insn
    val syscall : insn
	
    val alloc : label * int -> asm
    val alloc_float : label * real -> asm



    val print_program : program * TextIO.outstream -> unit
	
    val isPseudoReg : reg -> bool
    val isPseudoFReg : freg -> bool

  end (* signature MIPS *)
