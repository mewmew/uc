define void @f() {
; <label>:0
	%x = alloca i32
	%1 = load i32, i32* %x
	%2 = icmp ne i32 %1, 0
	br i1 %2, label %3, label %4

; <label>:3
	br label %5

; <label>:4
	br label %5

; <label>:5
	ret void
}
