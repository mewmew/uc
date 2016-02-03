/* Test file for semantic errors. Contains exactly one error. */

char a[10];

int main (void) {
  a = 42;	  // assign int to array of char
}
