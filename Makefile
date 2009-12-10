include $(GOROOT)/src/Make.$(GOARCH)

TARG=conc
GOFILES=\
        conc.go\
		filter.go\
		fold.go\
        for.go\
		future.go\
		map.go\
		realize.go\
		wait.go\

include $(GOROOT)/src/Make.pkg
