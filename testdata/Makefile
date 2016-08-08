all: uc.ll extra/irgen

uc.ll: uc.c
	clang -S -emit-llvm -o $@ $<
	./extra/irgen/strip.sh $@

extra/irgen:
	make -C $@

.PHONY: extra/irgen
