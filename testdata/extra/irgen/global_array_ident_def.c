// NOTE: Assignments to array references are not part of uC.
int b[5];
void f(int a[]) {
	a = b;
}
