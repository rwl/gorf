# Makefile generated by gb: http://go-gb.googlecode.com
# gb provides configuration-free building and distributing

include $(GOROOT)/src/Make.inc

TARG=gorf
GOFILES=\
	dscan.go\
	field.go\
	gorf.go\
	help.go\
	merge.go\
	move.go\
	pkg.go\
	rename.go\
	scan.go\
	singlemover.go\
	source.go\
	undo.go\
	util.go\

include $(GOROOT)/src/Make.cmd

_go_.$O: .goinstall

.goinstall:
	goinstall gonicetrace.googlecode.com/hg/nicetrace && \
	goinstall rog-go.googlecode.com/hg/exp/go/types && \
	goinstall rog-go.googlecode.com/hg/exp/go/parser && \
	touch .goinstall