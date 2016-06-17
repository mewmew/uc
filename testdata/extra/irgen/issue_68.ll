define i32 @main() {
0:
	%1 = alloca i32
	%a = alloca i32
	store i32 0, i32* %1
	%2 = load i32, i32* %a
	%3 = icmp ne i32 %2, 0
	br i1 %3, label %4, label %9
4:
	%5 = load i32, i32* %a
	%6 = icmp ne i32 %5, 0
	br i1 %6, label %7, label %8
7:
	store i32 42, i32* %a
	br label %8
8:
	br label %9
9:
	%10 = load i32, i32* %1
	ret i32 %10
}
