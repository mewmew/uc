define i32 @f() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = icmp sge i32 %1, %2
	%4 = zext i1 %3 to i32
	ret i32 %4
}
