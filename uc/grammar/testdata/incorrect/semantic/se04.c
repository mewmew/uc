/* Test file for semantic errors. Contains exactly one error. */

int a;

char a;		// Redeclaration of 'a'
int main(void) {
  a = 0;
}
