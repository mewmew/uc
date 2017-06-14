(* main/make.sml *)

(*val _ = CM.autoload "smlnj/cm/full.cm";*)
val _ = CM.make "sources.cm";
(*val _ = CM.make();*)
(*val _ = CM.State.reset();*)
(*val _ = CM.clear();*)
val _ = SMLofNJ.exportFn("ucc", Start.start Main.main);
