void f(int a[]) {
	;
}

int main(void) {
	typedef int foo;
	foo x[20];
	f(x);
}
