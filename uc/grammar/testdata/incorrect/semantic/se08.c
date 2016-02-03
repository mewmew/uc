/* Test file for semantic errors. Contains exactly one error. */

int a(int n) {
  return;		// Void return from function
}
int main(void) {
 a(2);
}
