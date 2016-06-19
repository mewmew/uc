define i32 @f() {
; <label>:0
	%a = alloca i32
	%b = alloca i32
	%1 = load i32, i32* %b
	store i32 %1, i32* %a
	ret i32 %1
}
