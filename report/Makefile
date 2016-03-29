all: uc.pdf clean

uc.pdf: inc
	pdflatex uc.tex
	bibtex uc
	pdflatex uc.tex
	pdflatex uc.tex

# *.tex -> *.pdf
TEX = $(wildcard *.tex)
PDF = $(TEX:.tex=.pdf)

.PHONY: clean inc

inc:
	make -C $@/

clean:
	rm -f *.aux *.log *.out *.toc *.bbl *.blg *.dvi *.rtf *.bcf *-blx.bib *.run.xml sections/*.aux

cleanall: clean
	rm -f uc.pdf 1_intro.pdf 2_lex.pdf 3_parse.pdf 4_sem.pdf 5_ir.pdf 6_ir_gen.pdf 7_con.pdf app.pdf
	make -C inc/ clean

# TODO: Remove temporary sections from Makefile when the report is finalized.

%.pdf: %.tex
	pdflatex $<
	pdflatex $<
	pdflatex $<

2_lex.pdf: inc

tmpsec: 1_intro.pdf 2_lex.pdf 3_parse.pdf 4_sem.pdf 5_ir.pdf 6_ir_gen.pdf 7_con.pdf app.pdf
