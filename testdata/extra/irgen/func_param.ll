define i32 @f(i32 %a) {
; <label>:0
	%1 = alloca i32
	store i32 %a, i32* %1
	%2 = load i32, i32* %1
	ret i32 %2
}
