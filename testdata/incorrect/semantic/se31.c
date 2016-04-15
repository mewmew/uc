/* Test file for semantic errors. Contains exactly one error. */

int a;

void a(void); 		// Attempt to redefine  'a' as extern


int main(void) {
  42;
}
