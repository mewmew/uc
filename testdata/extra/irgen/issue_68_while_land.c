void f() {
	int a;
	int b;

	while (a && b) {
		a = 11;
		while (a && b) {
			a = 22;
		}
		a = 33;
	}
}
