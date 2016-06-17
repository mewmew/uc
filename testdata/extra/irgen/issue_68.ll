define void @f() {
0:
	%a = alloca i32
	%1 = load i32, i32* %a
	%2 = icmp ne i32 %1, 0
	br i1 %2, label %3, label %8
3:
	%4 = load i32, i32* %a
	%5 = icmp ne i32 %4, 0
	br i1 %5, label %6, label %7
6:
	store i32 42, i32* %a
	br label %7
7:
	br label %8
8:
	ret void
}
