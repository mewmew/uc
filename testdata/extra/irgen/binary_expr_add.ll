define i32 @f() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = add i32 %1, %2
	ret i32 %3
}
define i32 @g() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = sub i32 %1, %2
	ret i32 %3
}
define i32 @h() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = mul i32 %1, %2
	ret i32 %3
}
define i32 @i() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = sdiv i32 %1, %2
	ret i32 %3
}
define i32 @j() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = icmp slt i32 %1, %2
	%4 = zext i1 %3 to i32
	ret i32 %4
}
define i32 @k() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = icmp sgt i32 %1, %2
	%4 = zext i1 %3 to i32
	ret i32 %4
}
define i32 @l() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = icmp sle i32 %1, %2
	%4 = zext i1 %3 to i32
	ret i32 %4
}
define i32 @m() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = icmp sge i32 %1, %2
	%4 = zext i1 %3 to i32
	ret i32 %4
}
define i32 @n() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = icmp ne i32 %1, %2
	%4 = zext i1 %3 to i32
	ret i32 %4
}
define i32 @o() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = icmp eq i32 %1, %2
	%4 = zext i1 %3 to i32
	ret i32 %4
}
define i32 @p() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = icmp ne i32 %1, 0
	br i1 %2, label %3, label %6
3:
	%4 = load i32, i32* %b
	%5 = icmp ne i32 %4, 0
	br label %6
6:
	%7 = phi i1 [ false, %0 ], [ %5, %3 ]
	%8 = zext i1 %7 to i32
	ret i32 %8
}
define i32 @q() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %b
	store i32 %1, i32* %a
	ret i32 %1
}
