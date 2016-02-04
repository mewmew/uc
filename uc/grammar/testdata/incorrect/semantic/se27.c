/* Test file for semantic errors. Contains exactly one error. */

void a (int n) {
  if (1<2) return 2 * n;	// Attempt to return value from procedure
  else 42;
}

int main(void) {
  a(2);
}

