FROM golang:1.15-alpine as build
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build ./cmd/control

FROM alpine:3.12
WORKDIR /
RUN addgroup -S systemstat
RUN adduser -S systemstat -G systemstat
COPY --from=build /build/control /control
RUN chown systemstat:systemstat /control
RUN chown -R systemstat:systemstat /control
ENTRYPOINT ["/control"]
USER systemstat
