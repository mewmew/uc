define i32 @f() {
; <label>:0
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %a
	%2 = load i32, i32* %b
	%3 = sub i32 %1, %2
	ret i32 %3
}
