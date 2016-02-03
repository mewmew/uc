/* Test file for semantic errors. Contains exactly one error. */

int a(void) {
  a = 1;	// 'a' is not an lval
  return 0;
}
int main(void) {
 a();
}

