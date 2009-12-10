include $(GOROOT)/src/Make.$(GOARCH)

TARG=conc
GOFILES=\
        conc.go\
		filter.go\
        for.go\
		future.go\
		map.go\
		realize.go\
		reduce.go\

include $(GOROOT)/src/Make.pkg
