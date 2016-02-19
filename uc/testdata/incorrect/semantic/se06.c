/* Test file for semantic errors. Contains exactly one error. */

int a(int n) {
  return 2 * n;
}

int a(int i) {   // Redeclaration of 'a'
  return i / 2;
}

int main(void) {
  a(2);
}
