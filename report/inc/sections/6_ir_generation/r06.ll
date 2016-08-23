define i32 @fac(i32 %n) {
; <label>:0
	%1 = alloca i32
	store i32 %n, i32* %1
	%i = alloca i32
	%p = alloca i32
	%2 = load i32, i32* %1
	%3 = icmp slt i32 %2, 0
	br i1 %3, label %4, label %5
; <label>:4
	ret i32 0
; <label>:5
	store i32 0, i32* %i
	store i32 1, i32* %p
	br label %6
; <label>:6
	%7 = load i32, i32* %i
	%8 = load i32, i32* %1
	%9 = icmp slt i32 %7, %8
	br i1 %9, label %10, label %16
; <label>:10
	%11 = load i32, i32* %i
	%12 = add i32 %11, 1
	store i32 %12, i32* %i
	%13 = load i32, i32* %p
	%14 = load i32, i32* %i
	%15 = mul i32 %13, %14
	store i32 %15, i32* %p
	br label %6
; <label>:16
	%17 = load i32, i32* %p
	ret i32 %17
}
define i32 @main() {
; <label>:0
	%1 = call i32 @fac(i32 5)
	ret i32 %1
}
