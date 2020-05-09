define void @f() {
0:
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = icmp ne i32 %1, 0
	br i1 %2, label %3, label %15
3:
	%4 = load i32, i32* %b
	%5 = icmp ne i32 %4, 0
	br i1 %5, label %6, label %15
6:
	%7 = load i32, i32* %a
	%8 = icmp ne i32 %7, 0
	br i1 %8, label %9, label %13
9:
	%10 = load i32, i32* %b
	%11 = icmp ne i32 %10, 0
	br i1 %11, label %12, label %13
12:
	store i32 11, i32* %a
	br label %14
13:
	store i32 22, i32* %a
	br label %14
14:
	store i32 33, i32* %a
	br label %16
15:
	store i32 44, i32* %a
	br label %16
16:
	ret void
}
