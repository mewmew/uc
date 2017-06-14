(* config/start.sml *)

signature START =
  sig
    val start: (string list -> unit) -> string * string list -> OS.Process.status
  end (* signature START *)

structure Start : START =
  struct

    fun sayErr msg = TextIO.output(TextIO.stdErr, msg)

    exception Interrupt

    (* This function applies operation to argument. If it handles an interrupt
     * signal (Control-C), it raises the exception Interrupt. Example:
     * (handleInterrupt(foo,x)) handle Interrupt => print "Bang!\n"
     *)
    fun handleInterrupt(operation: 'a -> unit, argument: 'a) =
      let exception Done
	  val oldHandler = Signals.inqHandler(Signals.sigINT)
	  fun resetHandler() =
	    Signals.setHandler(Signals.sigINT, oldHandler)
      in (SMLofNJ.Cont.callcc(fn k =>
		(Signals.setHandler(Signals.sigINT, Signals.HANDLER(fn _ => k));
		 operation argument;
		 raise Done));
	  raise Interrupt)
	 handle Done => (resetHandler(); ())
	      | exn  => (resetHandler(); raise exn)
      end

    fun start main (arg0,argv) =
      (handleInterrupt(main, argv); OS.Process.success)
      handle Interrupt => (sayErr "Interrupt\n"; OS.Process.failure)
	   | exn => (sayErr "Error: "; sayErr(General.exnMessage exn);
		     sayErr "\n"; OS.Process.failure)

  end (* structure Start *)
