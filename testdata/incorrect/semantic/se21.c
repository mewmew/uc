/* Test file for semantic errors. Contains exactly one error. */

int a(int n) {
  char bv[10];
    return bv;  //  Return from function with erroneous type
}
int main (void) {
  a(2);
}
