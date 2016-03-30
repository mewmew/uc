/* Test file for semantic errors. Contains exactly one error. */



int first(int b) {
  return b[0]; //not an array!
}

int main (void) {
  first(42);		
}
