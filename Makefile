# Author:  Tan Menglong <tanmenglong@gmail.com>
# Date:    Mon Dec  8 12:25:53 2014
#
# Make Target:
# ------------
# The Makefile provides the following targets to make:
#   $ make           compile and link
#   $ make clean     clean objects and the executable file
#
#===========================================================================

.PHONY : all dev deps output clean help test deb rpm doc

all : output

dev : deps
	go get github.com/goccmack/gocc

deps :
	go get github.com/colinmarc/hdfs
	go get github.com/awalterschulze/gographviz
	go get github.com/mattn/go-sqlite3
	go get github.com/crackcell/gotabulate

output : hpipe
	mkdir -p output/bin
	cp hpipe output/bin/hpipe

deb : output
	mkdir -p dist/deb/usr/bin
	cp output/bin/* dist/deb/usr/bin/
	mkdir -p output/packages
	dpkg -b ${PWD}/dist/deb output/packages/hpipe_`grep Version dist/deb/DEBIAN/control | cut -f 2 -d " "`_amd64.deb

rpm : output
	cd dist/rpm; ./hpipe-build.sh

hpipe :
	go build

test : all
	go test -v

doc :
	cd docs; make html

clean :
	rm -rf output
	rm -rf hpipe

help :
	@echo 'Usage: make [TARGET]'
	@echo 'TARGETS:'
	@echo '  all       (=make) compile and link.'
	@echo '  dev       get dependencies for development.'
	@echo '  deps      get dependencies for compiling.'
	@echo '  clean     clean objects and the executable file.'
	@echo '  help      print this message.'
