/* Test file for semantic errors. Contains exactly one error. */

void a (int n) {
  char n;
  42;
}

int main(void) {
  a(2);
}

