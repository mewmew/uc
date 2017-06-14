(* util/aa-dict-fn.sml
 * Written by: Mikael Pettersson, mpe@ida.liu.se, 1996.
 *
 * An SML implementation of dictionaries based on balanced search
 * trees as described in the following article:
 *
 * @InProceedings{Andersson93,
 *  author =	{Arne Andersson},
 *  title =	{Balanced Search Trees Made Simple},
 *  pages =	{60--71},
 *  crossref =	{WADS93}
 * }
 *
 * @Proceedings{WADS93,
 *   title =	{Proceedings of the Third Workshop on Algorithms and Data Structures, WADS'93},
 *   booktitle ={Proceedings of the Third Workshop on Algorithms and Data Structures, WADS'93},
 *   editor =	{F. Dehne and J-R. Sack and N. Santoro and S. Whitesides},
 *   publisher ={Springer-Verlag},
 *   series =	{Lecture Notes In Computer Science},
 *   volume =	709,
 *   year =	1993
 * }
 *
 * The original Pascal code represented empty trees by a single shared node
 * with level 0, and whose left and right fields pointed back to itself.
 * The point of this trick was to eliminate special cases in the skew and
 * split procedures. Since it would be expensive to emulate, this SML code
 * uses a traditional representation, making all special cases explicit.
 *
 * This is the vanilla version with no optimizations applied.
 *)
functor AADictFn(Key : ORD_KEY) : ORD_DICT =
  struct

    structure Key = Key

    local

      datatype 'a tree	= E
			| T of {key: Key.ord_key,
				attr: 'a,
				level: int,
				left: 'a tree,
				right: 'a tree	}

      fun split(t as E) = t
	| split(t as T{right=E,...}) = t
	| split(t as T{right=T{right=E,...},...}) = t
	| split(t as T{key=kx,attr=ax,level=lx,left=a,
		       right=T{key=ky,attr=ay,left=b,
			       right=(z as T{level=lz,...}),...}}) =
	    if lx = lz then	(* rotate left *)
	      T{key=ky,attr=ay,level=lx+1,right=z,
		left=T{key=kx,attr=ax,level=lx,left=a,right=b}}
	    else t

      fun skew(t as E) = t
	| skew(t as T{left=E,...}) = t
	| skew(t as T{key=kx,attr=ax,level=lx,right=c,
		      left=T{key=ky,attr=ay,level=ly,left=a,right=b}}) =
	    if lx = ly then	(* rotate right *)
	      T{key=ky,attr=ay,level=ly,left=a,
		right=T{key=kx,attr=ax,level=lx,left=b,right=c}}
	    else t

      fun tfind(t, x) =
	let fun look(E) = E
	      | look(t as T{key,left,right,...}) =
		  case Key.compare(x, key)
		    of LESS => look left
		     | GREATER => look right
		     | EQUAL => t
	in
	  look t
	end

    in

      type 'a dict = 'a tree
      val empty = E

      fun insert(E, x, y) = T{key=x, attr=y, level=1, left=E, right=E}
	| insert(T{key,attr,level,left,right}, x, y) =
	    let val t = case Key.compare(x,key)
			  of LESS =>
			      T{key=key, attr=attr, level=level, right=right,
				left=insert(left,x,y)}
			   | GREATER =>
			      T{key=key, attr=attr, level=level, left=left,
				right=insert(right,x,y)}
			   | EQUAL =>
			      T{key=x, attr=y, level=level, left=left, right=right}
		val t = skew t
		val t = split t
	    in
	      t
	    end

      fun find'(t, x) =
	case tfind(t, x)
	  of E => NONE
	   | T{key,attr,...} => SOME(key,attr)

      fun find(t, x) =
	case tfind(t, x)
	  of E => NONE
	   | T{attr,...} => SOME attr

      fun plus(bot, E) = bot
	| plus(bot, T{key,attr,left,right,...}) =
	    insert(plus(plus(bot, left), right), key, attr)

      fun fold(f, init, dict) =
	let fun traverse(E, state) = state
	      | traverse(T{key,attr,left,right,...}, state) =
		  traverse(right, traverse(left, f(key,attr,state)))
	in
	  traverse(dict, init)
	end

    end (* local *)

  end (* functor AADictFn *)
