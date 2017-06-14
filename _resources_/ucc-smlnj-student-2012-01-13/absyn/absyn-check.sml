(* absyn/absyn-check.sml *)

signature ABSYN_CHECK =
  sig
    structure Absyn: ABSYN
    val program: Absyn.program -> unit
  end (* signature ABSYN_CHECK *)

functor AbsynCheckFn(Absyn : ABSYN) : ABSYN_CHECK =
  struct

    structure Absyn = Absyn

    (*
     * Reporting errors.
     *
     * Source file context is not easily available everywhere, so
     * a detected error is instead thrown as an exception.
     * At the top level where we do have the source file context,
     * we catch this exception and generate appropriate messages
     * before re-throwing the exception.
     * Limitation: We can't continue after an error. Big deal.
     *)

    type msg = string * int * int (* same as what Absyn.Source.sayMsg wants *)
    exception AbsynCheckError of msg list

    fun withSource(source, f) =
      f()
      handle (exn as AbsynCheckError(msgs)) =>
	(List.app (Absyn.Source.sayMsg source) msgs;
	 raise exn)

    fun error1 msg = raise AbsynCheckError[msg]

    fun error2(msg1, msg2) = raise AbsynCheckError[msg1, msg2]

    fun mkIdErrorMsg(msg, Absyn.IDENT(name, left, right)) =
      ("Error: "^msg^name, left, right)
    fun idError(msg, id) = error1(mkIdErrorMsg(msg, id))

    fun doError(msg, left, right) = error1("Error: "^msg, left, right)
    fun expError(msg, Absyn.EXP(_,left,right)) = doError(msg, left, right)
    fun stmtError(msg, Absyn.STMT(_,left,right)) = doError(msg, left, right)

    (*
     * YOUR CODE HERE
     *
     * Hints:
     * - You need to represent uC types.
     * - You need an environment/symbol-table for identifiers.
     * - You need recursive functions over expressions and statements.
     * - You need to check type constraints at various places.
     * - Abstract syntax 'declarators' aren't types. You'll need
     *   to translate them.
     * - You need to process top-level declarations.
     *)

    fun checkDeclarations _ = ()	(* XXX: REPLACE WITH YOUR CODE *)

    (* Programs *)

    fun program(Absyn.PROGRAM{decs,source}) =
      let fun check() = checkDeclarations decs
      in
	withSource(source, check)
      end

  end (* functor AbsynCheckFn *)
