define i32 @f() {
0:
	%a = alloca i32
	%1 = load i32, i32* %a
	%2 = icmp ne i32 %1, 0
	%3 = xor i1 %2, true
	%4 = zext i1 %3 to i32
	ret i32 %4
}
