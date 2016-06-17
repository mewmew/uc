void f(){
	int a;

	if (a) {
		a = 11;
		while (a) {
			a = 22;
			if (a) {
				a = 33;
				while (a) {
					a = 44;
				}
				a = 55;
			} else {
				a = 66;
				while (a) {
					a = 77;
				}
				a = 88;
			}
			a = 99;
		}
		a = 111;
	} else {
		a = 222;
	}
}
