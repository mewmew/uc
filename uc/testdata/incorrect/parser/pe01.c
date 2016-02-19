/* Test file for syntactic errors. Contains exactly one error. */

int a;
int main(void) {
  a = (a + ) * a;			//	 Unexpected token ')'
}
