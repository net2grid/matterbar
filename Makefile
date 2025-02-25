.PHONY: build

BUILD_TAG = $(shell git describe)

build:
	mkdir -p build
	cd server && env GOOS=linux GOARCH=amd64 go build -o ../build/plugin-linux-amd64;
	cd server && env GOOS=darwin GOARCH=amd64 go build -o ../build/plugin-darwin-amd64;
	cd server && env GOOS=windows GOARCH=amd64 go build -o ../build/plugin-windows-amd64.exe;

build-audit:
	cd scripts && go build -o audit

bundle:
	rm -rf dist/matterbar/
	mkdir -p dist/matterbar/server/
	cp plugin.json dist/matterbar/
	cp -r assets dist/matterbar/
	cp -r build/* dist/matterbar/server/
	cd dist && tar -cvzf "matterbar-$(BUILD_TAG).tar.gz" matterbar/

clean:
	rm -rf dist/

coverage: test
	go tool cover -html=server/coverage.txt

deploy: build bundle
ifneq ($(wildcard ../mattermost-server/.*),)
	mkdir -p ../mattermost-server/plugins
	tar -C ../mattermost-server/plugins -zxvf "dist/matterbar-$(BUILD_TAG).tar.gz"
else
	@echo "Unable to find local mattermost-server dir. Try installing manually."
endif

deps:
	dep ensure

deps-update:
	dep ensure -update

fmt:
	gofmt -w server scripts

test:
	cd server && go test -v -race -coverprofile coverage.txt
