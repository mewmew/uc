define void @f() {
; <label>:0
	%x = alloca i32
	br label %1

; <label>:1
	%2 = load i32, i32* %x
	%3 = icmp ne i32 %2, 0
	br i1 %3, label %4, label %5

; <label>:4
	br label %1

; <label>:5
	ret void
}
