default: deps mybuild

deps:
	go get

mybuild:
	go build rcadmin.go

xcompile:
	goreleaser --snapshot --skip-publish --rm-dist

clean:
	rm -f rcadmin
	rm -rf dist/*
