PWD=$(shell pwd)
PKGNAME=activityjs.io/serve
EXCLUDE_OPTS=--exclude-from=.gitignore --exclude=.git --exclude=go.mod  --exclude=workspace
DEBUGOPT=--tags=debug

.PHONY: demo path

build:
	- mkdir -p src/activityjs.io
	- rm -f src/activityjs.io/serve
	- ln -s $(PWD) src/activityjs.io/serve

demo:
	@- GOPATH=$(PWD) gopherjs build demo/main.go -o demo/main.js
	@ls -lah demo

debug:
	@- GOPATH=$(PWD) gopherjs build demo/main.go -v --tags="debug" -o demo/main.js
	@ls -lah demo

serve:
	@go run demo/serve.go

path:
	- GOPATH=$(PWD) $(run)

makeworker:
	@-mkdir -p workspace/src/$(PKGNAME)
	@-mkdir workspace/tmp
	@rsync -rv $(EXCLUDE_OPTS) . workspace/src/$(PKGNAME)
	@./scripts/install.sh $(PWD)/workspace

