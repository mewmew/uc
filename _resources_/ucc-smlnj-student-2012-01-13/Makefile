# Makefile for the Micro-C compiler
SHELL=/bin/sh

GENERATED=lexer/uc.lex.sml parser/uc.grm.sml parser/uc.grm.sig parser/uc.grm.desc

all:
	cat main/make.sml | sml

distclean realclean:	clean
	rm -rf */.cm
	rm -f $(GENERATED)

clean:
	rm -f ucc.*-*
