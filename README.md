# systemstat

work in progress

## Developing

Systemstat is developed with Golang 1.15.

Some services leverage GRPC and protocol buffers
- install `protoc` from [Github](https://github.com/protocolbuffers/protobuf/releases)
```
# Linux Example
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip -O protoc.zip
unzip -j protoc.zip bin/protoc -d ./
sudo mv protoc /usr/local/bin/protoc
rm -f protoc.zip
```
- install required go modules:
```
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u google.golang.org/grpc
```

## Testing

Make sure `systemstat` is [in your Gopath](https://golang.org/doc/code.html)
- `make test`: shellcheck and go tests
- `make test.integration`: end to end testing (relies on shellcheck, docker, docker-compose, and some local horsepower)

## Faq

Q: Why is Systemstat a monorepo?
A: Because this project is new and unstable. Once the structure and APIs are well defined, services may move to their own repo.
