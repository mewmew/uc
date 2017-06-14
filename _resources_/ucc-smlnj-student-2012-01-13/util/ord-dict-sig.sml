(* util/ord-dict-sig.sml *)

signature ORD_DICT =
  sig

    structure Key : ORD_KEY

    type 'a dict
    val empty	: 'a dict
    val insert	: 'a dict * Key.ord_key * 'a -> 'a dict
    val find'	: 'a dict * Key.ord_key -> (Key.ord_key * 'a) option
    val find	: 'a dict * Key.ord_key -> 'a option
    val plus	: 'a dict * 'a dict -> 'a dict
    val fold	: (Key.ord_key * 'a * 'b -> 'b) * 'b * 'a dict -> 'b

  end (* signature ORD_DICT *)
