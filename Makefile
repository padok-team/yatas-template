default: build

test:
	go test ./...

build:
	go build -o bin/yatas-template

update:
	go get -u
	go mod tidy

install: build
	mkdir -p ~/.yatas.d/plugins/github.com/padok-team/yatas-template/local/
	mv ./bin/yatas-template ~/.yatas.d/plugins/github.com/padok-team/yatas-template/local/

release: test
	npx standard-version
	git push --follow-tags origin main