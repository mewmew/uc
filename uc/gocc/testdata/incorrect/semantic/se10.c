/* Test file for semantic errors. Contains exactly one error. */

int n;
int main(void) {
  n = 42;
  n[2]; // Index an integer
}
