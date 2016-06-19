define i32 @f() {
; <label>:0
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = icmp ne i32 %1, 0
	br i1 %2, label %3, label %6
; <label>:3
	%4 = load i32, i32* %b
	%5 = icmp ne i32 %4, 0
	br label %6
; <label>:6
	%7 = phi i1 [ false, %0 ], [ %5, %3 ]
	%8 = zext i1 %7 to i32
	ret i32 %8
}
