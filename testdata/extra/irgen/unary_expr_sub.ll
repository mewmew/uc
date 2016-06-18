define i32 @f() {
; <label>:0
	%a = alloca i32
	%1 = load i32, i32* %a
	%2 = sub i32 0, %1
	ret i32 %2
}
