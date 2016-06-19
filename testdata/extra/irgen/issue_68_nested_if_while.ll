define void @f() {
; <label>:0
	%a = alloca i32
	%1 = load i32, i32* %a
	%2 = icmp ne i32 %1, 0
	br i1 %2, label %3, label %24
; <label>:3
	store i32 11, i32* %a
	br label %4
; <label>:4
	%5 = load i32, i32* %a
	%6 = icmp ne i32 %5, 0
	br i1 %6, label %7, label %23
; <label>:7
	store i32 22, i32* %a
	%8 = load i32, i32* %a
	%9 = icmp ne i32 %8, 0
	br i1 %9, label %10, label %16
; <label>:10
	store i32 33, i32* %a
	br label %11
; <label>:11
	%12 = load i32, i32* %a
	%13 = icmp ne i32 %12, 0
	br i1 %13, label %14, label %15
; <label>:14
	store i32 44, i32* %a
	br label %11
; <label>:15
	store i32 55, i32* %a
	br label %22
; <label>:16
	store i32 66, i32* %a
	br label %17
; <label>:17
	%18 = load i32, i32* %a
	%19 = icmp ne i32 %18, 0
	br i1 %19, label %20, label %21
; <label>:20
	store i32 77, i32* %a
	br label %17
; <label>:21
	store i32 88, i32* %a
	br label %22
; <label>:22
	store i32 99, i32* %a
	br label %4
; <label>:23
	store i32 111, i32* %a
	br label %25
; <label>:24
	store i32 222, i32* %a
	br label %25
; <label>:25
	ret void
}
