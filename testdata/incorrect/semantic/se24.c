/* Test file for semantic errors. Contains exactly one error. */

char first(int a[]) {
  int b[10];

  b = a;		// b cannot be assigned
  return 1;
}

