
install:
	go install

lint:
	gometalinter --config .linter.conf

test:
	go test --cover -v $$(go list ./... | grep -v /vendor/)

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

dev-deps: deps
	go get github.com/smartystreets/goconvey
	go get github.com/alecthomas/gometalinter
	gometalinter --install
