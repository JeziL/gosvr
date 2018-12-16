gosvr: static
	go build

install: gosvr
	go install

ds:
	find . -name '*.DS_Store' -type f -delete

static: ds
	packr2

clean: ds
	rm -f gosvr
	rm -rf release

cleanstatic: ds
	rm -rf packrd
	rm -f *packr.go

cleanall: clean cleanstatic
