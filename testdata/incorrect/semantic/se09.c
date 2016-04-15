/* Test file for semantic errors. Contains exactly one error. */

int a(int n) {
  char x[1];
  if (n != 0) return 2;
  else return x;    // Return from function with erroneous type
}

int main(void) {
  a(5);
}


