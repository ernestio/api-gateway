install:
	go install

lint:
	# golint ./...
	# go vet ./...

test:
	go test -v ./... --cover

deps: 
	go get github.com/Masterminds/glide
	glide install

dev-deps: deps
