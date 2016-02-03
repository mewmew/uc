/* Test file for syntactic errors. Contains exactly one error. */

int bar(void) {
  int foo(void) {	// Local procedure definitions are not allowed
    ;
  }
  ;
}



