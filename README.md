# systemstat

work in progress

## Developing

Systemstat is developed with Golang 1.15.

## Testing

Make sure `systemstat` is [in your Gopath](https://golang.org/doc/code.html)
- `make test`: shellcheck and go tests
- `make test.integration`: end to end testing (relies on shellcheck, docker, docker-compose, and some local horsepower)

## Faq

Q: Why is Systemstat a monorepo?
A: Because this project is new and unstable. Once the structure and APIs are well defined, services may move to their own repo.
