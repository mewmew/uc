define i32 @f() {
; <label>:0
	%1 = alloca i32
	%a = alloca i32
	br label %2
; <label>:2
	%3 = load i32, i32* %a
	%4 = icmp ne i32 %3, 0
	br i1 %4, label %5, label %6
; <label>:5
	store i32 1, i32* %1
	br label %7
; <label>:6
	store i32 2, i32* %1
	br label %7
; <label>:7
	%8 = load i32, i32* %1
	ret i32 %8
}
