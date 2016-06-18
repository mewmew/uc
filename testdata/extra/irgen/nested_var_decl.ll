@a = global i32 0
define i32 @f() {
; <label>:0
	%a = alloca i32
	%a1 = alloca i32
	%1 = load i32, i32* @a
	ret i32 %1
}
