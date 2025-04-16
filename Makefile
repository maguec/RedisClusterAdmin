default: deps mybuild

deps:
	go get

mybuild:
	go build -o RedisClusterAdmin rcadmin.go 

xcompile:
	goreleaser --snapshot --skip-publish --rm-dist

clean:
	rm -f RedisClusterAdmin
	rm -rf dist/*
