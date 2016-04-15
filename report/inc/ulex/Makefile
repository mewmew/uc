.SECONDARY:
TESTDATA=../../../testdata
VPATH=$(CURDIR):$(TESTDATA)
UPROG=ulex

all: clean lexical_analysis

LEXSOURCEDIRS=$(wildcard $(TESTDATA)/*/lexer)
LEXSOURCE=$(wildcard $(LEXSOURCEDIRS:=/*.c))
LEXOBJECTDIRS=$(LEXSOURCEDIRS:$(TESTDATA)/%=%)
LEXOBJECTS=$(LEXSOURCE:$(TESTDATA)/%.c=%.txt)

UCFLAGS=
REDIRECTS=2>&1

lexical_analysis: $(LEXOBJECTDIRS:=/listing.tex)

$(LEXOBJECTDIRS):
	mkdir -p "$@"

$(LEXOBJECTS): %.txt: $(TESTDATA)/%.c | $(LEXOBJECTDIRS)
	cd $(CURDIR)/$(TESTDATA); $(UPROG) $(UCFLAGS) $(<:$(TESTDATA)/%=%) > $(CURDIR:$(TESTDATA)/%=%)/$@ $(REDIRECTS)

%/listing.tex: $(LEXOBJECTS)
	ls $(@D)/*.txt | awk '{name=$$1;\
		sub(/.txt/,".c",name); printf "\\lstinputlisting[style=go,\
		caption=%s,label=lst:$(UPROG)/%s]{inc/$(UPROG)/%s}\n\n", name, name, $$1}' > $@
	ls $(@D)/*.txt | awk '{name=$$1;\
		sub(/.txt/,".c",name); printf "\nListing~\\ref{lst:$(UPROG)/%s}:~\\nameref{lst:$(UPROG)/%s}\n",\
		name, name}' > $(@D)/listref.tex

incorrect/lexer/ugly.txt: UCFLAGS+=-n 10

incorrect/lexer/%.txt: REDIRECTS+=1>/dev/null

clean:
	rm -rf $(LEXOBJECTDIRS)

.PHONY: all lexical_analysis $(LEXOBJECTDIRS) clean 
