define void @f() {
; <label>:0
	%a = alloca i32
	%1 = load i32, i32* %a
	ret void
}
