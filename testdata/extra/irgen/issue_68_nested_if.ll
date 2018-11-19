define void @f() {
; <label>:0
	%a = alloca i32
	%1 = load i32, i32* %a
	%2 = icmp ne i32 %1, 0
	br i1 %2, label %3, label %9

; <label>:3
	%4 = load i32, i32* %a
	%5 = icmp ne i32 %4, 0
	br i1 %5, label %6, label %7

; <label>:6
	store i32 11, i32* %a
	br label %8

; <label>:7
	store i32 22, i32* %a
	br label %8

; <label>:8
	store i32 33, i32* %a
	br label %10

; <label>:9
	store i32 44, i32* %a
	br label %10

; <label>:10
	ret void
}
