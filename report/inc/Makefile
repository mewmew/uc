all: 2_lexical_analysis

2_lexical_analysis:
	ulex -n 10 ../../uc/testdata/quiet/lexer/*.c > 2_lexical_analysis/3_validation/quiet.txt 2>&1
	ulex -n 10 ../../uc/testdata/incorrect/lexer/*.c > 2_lexical_analysis/3_validation/incorrect.txt 2>&1

.PHONY: 2_lexical_analysis
