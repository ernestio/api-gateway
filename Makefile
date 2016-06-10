install:
	go install

deps:
	go get github.com/nats-io/nats
	go get github.com/labstack/echo
	go get github.com/dgrijalva/jwt-go
