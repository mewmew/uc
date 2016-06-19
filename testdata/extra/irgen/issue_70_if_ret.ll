define i32 @f() {
; <label>:0
	%1 = alloca i32
	%a = alloca i32
	%2 = load i32, i32* %a
	%3 = icmp ne i32 %2, 0
	br i1 %3, label %4, label %5
; <label>:4
	store i32 1, i32* %1
	br label %6
; <label>:5
	store i32 2, i32* %1
	br label %6
; <label>:6
	%7 = load i32, i32* %1
	ret i32 %7
}
