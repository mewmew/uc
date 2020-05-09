define void @f() {
0:
	%a = alloca i32
	%b = alloca i32
	br label %1

1:
	%2 = load i32, i32* %a
	%3 = icmp ne i32 %2, 0
	br i1 %3, label %4, label %7

4:
	%5 = load i32, i32* %b
	%6 = icmp ne i32 %5, 0
	br label %7

7:
	%8 = phi i1 [ false, %1 ], [ %6, %4 ]
	br i1 %8, label %9, label %20

9:
	store i32 11, i32* %a
	br label %10

10:
	%11 = load i32, i32* %a
	%12 = icmp ne i32 %11, 0
	br i1 %12, label %13, label %16

13:
	%14 = load i32, i32* %b
	%15 = icmp ne i32 %14, 0
	br label %16

16:
	%17 = phi i1 [ false, %10 ], [ %15, %13 ]
	br i1 %17, label %18, label %19

18:
	store i32 22, i32* %a
	br label %10

19:
	store i32 33, i32* %a
	br label %1

20:
	ret void
}
