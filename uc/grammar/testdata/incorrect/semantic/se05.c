/* Test file for semantic errors. Contains exactly one error. */

int a;

void a(void) {		// Attempt to redefine variable 'a'
  0;
}

int main(void) {
  42;
}
