.SECONDARY:
TESTDATA=../../../testdata
VPATH=$(TESTDATA)
UPROG=uparse

all: clean syntactic_analysis

SOURCES=quiet/lexer/l05.c quiet/parser/p01.c quiet/parser/p02.c quiet/parser/p03.c quiet/parser/p04.c quiet/parser/p07.c quiet/parser/p08.c incorrect/parser/pe01.c incorrect/parser/pe03.c incorrect/parser/pe05.c incorrect/parser/pe10.c incorrect/parser/pe11.c incorrect/parser/pe12.c
PARSESOURCEDIRS=$(sort $(dir $(SOURCES)))
PARSESOURCE=$(addprefix $(TESTDATA)/,$(SOURCES))
PARSEOBJECTDIRS=$(PARSESOURCEDIRS:$(TESTDATA)/%=%)
PARSEOBJECTS=$(PARSESOURCE:$(TESTDATA)/%.c=%.txt)

UCFLAGS=
REDIRECTS=2>&1

$(PARSEOBJECTDIRS):
	-mkdir -p $@

$(PARSEOBJECTS): %.txt: $(TESTDATA)/%.c | $(PARSEOBJECTDIRS)
	cd $(CURDIR)/$(TESTDATA); $(UPROG) $(UCFLAGS) $(<:$(TESTDATA)/%=%) > $(CURDIR:$(TESTDATA)/%=%)/$@ $(REDIRECTS)

syntactic_analysis: $(PARSEOBJECTDIRS:%=%listing.tex)

%/listing.tex: $(PARSEOBJECTS)
	ls $(@D)/*.txt | awk '{name=$$1;\
		sub(/.txt/,".c",name); printf "\\lstinputlisting[style=go,\
		caption=%s,label=lst:$(UPROG)/%s]{inc/$(UPROG)/%s}\n\n", name, name, $$1}' > $@
	ls $(@D)/*.txt | awk '{name=$$1;\
		sub(/.txt/,".c",name); printf "\nListing~\\ref{lst:$(UPROG)/%s}:~\\nameref{lst:$(UPROG)/%s}\n",\
		name, name}' > $(@D)/listref.tex

clean:
	rm -rf $(PARSEOBJECTDIRS)

.PHONY: all syntactic_analysis $(PARSEOBJECTDIRS) clean
