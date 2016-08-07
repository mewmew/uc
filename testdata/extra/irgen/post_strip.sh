#!/bin/bash
for f in "array_arg.ll" "array_ident_use.ll" "array_ident_def.ll" "global_array_arg.ll" "global_array_ident_use.ll" "global_array_param.ll" "global_array_ident_def.ll"; do
	sar -i "getelementptr ([^,]+), ([^,]+), i32 0, i32 0" "getelementptr \$1, \$2, i64 0, i64 0" "${f}"
done
