.SECONDARY:
TESTDATA=../../testdata
TESTOUTDIR=testoutput
VPATH=$(TESTDATA):$(TESTOUTDIR)

all: clean lexical_analysis syntactic_analysis

LEXSOURCEDIRS=$(wildcard $(TESTDATA)/*/lexer)
LEXSOURCE=$(wildcard $(LEXSOURCEDIRS:=/*.c))
LEXOBJECTDIRS=$(LEXSOURCEDIRS:$(TESTDATA)/%=$(TESTOUTDIR)/%)
LEXOBJECTS=$(LEXSOURCE:$(TESTDATA)/%.c=$(TESTOUTDIR)/%.txt)

PARSESOURCEDIRS=$(wildcard $(TESTDATA)/*/parser)
PARSESOURCETMP=quiet/lexer/l05.c quiet/parser/p01.c quiet/parser/p02.c quiet/parser/p03.c quiet/parser/p04.c quiet/parser/p07.c quiet/parser/p08.c incorrect/parser/pe01.c incorrect/parser/pe03.c incorrect/parser/pe05.c incorrect/parser/pe10.c incorrect/parser/pe11.c incorrect/parser/pe12.c
PARSESOURCE=$(addprefix $(TESTDATA)/,$(PARSESOURCETMP))
PARSEOBJECTDIRS=$(PARSESOURCEDIRS:$(TESTDATA)/%=$(TESTOUTDIR)/%)
PARSEOBJECTS=$(PARSESOURCE:$(TESTDATA)/%.c=$(TESTOUTDIR)/%.txt)

UCFLAGS=
REDIRECTS=2>&1

$(LEXOBJECTDIRS) $(PARSEOBJECTDIRS) :
	-mkdir -p $@

$(LEXOBJECTS): | $(LEXOBJECTDIRS)

$(PARSEOBJECTS): | $(PARSEOBJECTDIRS)

lexical_analysis: $(LEXOBJECTDIRS:%=%/listing.tex)

syntactic_analysis: $(PARSEOBJECTDIRS:%=%/listing.tex)

%/listing.tex: $(LEXOBJECTS) $(PARSEOBJECTS) #%/*.txt  #$(filter %,$(LEXOBJECTS))
	ls $(@D)/*.txt | awk '{name=$$1;sub(/$(TESTOUTDIR)\//,"",name);\
		sub(/.txt/,".c",name); printf "\\lstinputlisting[style=go,\
		caption=%s,label=lst:%s]{inc/%s}\n\n", name, name, $$1}' > $@
	ls $(@D)/*.txt | awk '{name=$$1;sub(/$(TESTOUTDIR)\//,"",name);\
		sub(/.txt/,".c",name); printf "\nListing~\\ref{lst:%s}:~\\nameref{lst:%s}\n",\
		name, name}' > $(@D)/listref.tex

$(TESTOUTDIR)/incorrect/lexer/ugly.txt: UCFLAGS+=-n 10

$(TESTOUTDIR)/incorrect/lexer/%.txt: REDIRECTS+=1>/dev/null

$(PARSEOBJECTS): UPROG=uparse

$(LEXOBJECTS): UPROG=ulex

$(TESTOUTDIR)/%.txt: $(TESTDATA)/%.c | $(LEXOBJECTDIRS)
	cd $(CURDIR)/$(TESTDATA); $(UPROG) $(UCFLAGS) $(<:$(TESTDATA)/%=%) > $(CURDIR:$(TESTDATA)/%=%)/$@ $(REDIRECTS)

clean:
	rm -rf $(TESTOUTDIR)

.PHONY: all lexical_analysis syntactic_analysis clean $(LEXOBJECTDIRS)
