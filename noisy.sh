#!/bin/bash
for f in testdata/noisy/simple/*.c; do
	echo "UC"
	./compile.sh "${f}"
	echo -e "\n### UC output:"
	lli out.ll
	echo -e "\n\nClang"
	clang -o a.out "${f}" testdata/uc.ll

	echo -e "\n###Clang output:"
	./a.out
	echo -e "\n"
done
