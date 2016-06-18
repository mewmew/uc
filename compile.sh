#!/bin/bash
if [ $# -lt 1 ]; then
	echo "Usage: strip FILE.ll"
	exit 1
fi
f=$1
uclang -o a.ll "${f}"
llvm-link -S -o out.ll a.ll testdata/uc.ll

echo "SUCCESS: ${f}"
