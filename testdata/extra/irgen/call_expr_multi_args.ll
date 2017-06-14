define i32 @g(i32 %a, i32 %b) {
; <label>:0
	%1 = alloca i32
	%2 = alloca i32
	store i32 %a, i32* %1
	store i32 %b, i32* %2
	ret i32 42
}

define i32 @f() {
; <label>:0
	%x = alloca i32
	%y = alloca i32
	%1 = load i32, i32* %y
	%2 = load i32, i32* %x
	%3 = call i32 @g(i32 %1, i32 %2)
	ret i32 %3
}
