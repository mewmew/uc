/* Test file for semantic errors. Contains exactly one error. */

char a[10];
int main (void) {
  if (a==42) ;
   // Array of char in expression
}
