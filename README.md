# API Gateway

master:  [![CircleCI](https://circleci.com/gh/ernestio/api-gateway/tree/master.svg?style=shield)](https://circleci.com/gh/ernestio/api-gateway/tree/master)  
develop: [![CircleCI](https://circleci.com/gh/ernestio/api-gateway/tree/develop.svg?style=shield)](https://circleci.com/gh/ernestio/api-gateway/tree/develop)

This service expose REST endpoints to each of the datastores inside of ernest. Requests are translated from http calls to nats requests.

## Installation

```
make deps
make install
```

## Running Tests

```
make deps
go test
```

## Authentication

Authentication is handled by JWT. You must first authenticate via `/auth/` and use the returned web token as a header in all subsequent requests.

```
curl -i -X POST -d "username=something" -d "password=something" localhost:8080/auth/
```

This will return the following json payload:

```json
{"token":"VALID-AUTH-TOKEN"}
```

This then can be used in subsequent requests, like so:

```
curl -i -H 'Authorization: Bearer VALID-AUTH-TOKEN' localhost:8080/api/users/
```

## Endpoints

Supported endpoints are Users, Groups, Datacenters and Services.


## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2015 ernest.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).
