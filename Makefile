gosvr: static
	go build

install: gosvr
	go install

ds:
	find . -name '*.DS_Store' -type f -delete

static: ds
	packr

clean: ds
	rm -f gosvr

cleanstatic: ds
	rm -f *packr.go

cleanall: clean cleanstatic
