all: uc.pdf clean_aux

uc.pdf: inc
	pdflatex uc.tex
	bibtex uc
	pdflatex uc.tex
	pdflatex uc.tex

# *.tex -> *.pdf
TEX = $(wildcard *.tex)
PDF = $(TEX:.tex=.pdf)

.PHONY: clean inc inc/uparse inc/ulex inc/usem
inc inc/uparse inc/ulex inc/usem:
	make -C $@/

clean_aux:
	rm -f *.log *.out *.toc *.bbl *.blg *.dvi *.rtf *.bcf *-blx.bib *.run.xml
	find . -iname "*.aux" -type f -delete

clean: clean_aux
	rm -f uc.pdf 1_intro.pdf 2_lex.pdf 3_parse.pdf 4_sem.pdf 5_ir.pdf 6_ir_gen.pdf 7_con.pdf app.pdf
	make -C inc/ clean

# TODO: Remove temporary sections from Makefile when the report is finalized.

lex: 2_lex.pdf clean_aux
parse: 3_parse.pdf clean_aux
sem: 4_sem.pdf clean_aux
ir: 5_ir.pdf clean_aux

%.pdf: %.tex
	pdflatex $<
	bibtex $*
	pdflatex $<
	pdflatex $<

2_lex.pdf: inc/ulex

3_parse.pdf: inc/uparse

4_sem.pdf: inc/usem

tmpsec: 1_intro.pdf 2_lex.pdf 3_parse.pdf 4_sem.pdf 5_ir.pdf 6_ir_gen.pdf 7_con.pdf app.pdf
