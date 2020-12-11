FROM golang:1.15-alpine as build
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build ./cmd/alert

FROM alpine:3.12
WORKDIR /
RUN addgroup -S systemstat
RUN adduser -S systemstat -G systemstat
COPY --from=build /build/alert /alert
RUN chown systemstat:systemstat /alert
RUN chown -R systemstat:systemstat /alert
ENTRYPOINT ["/alert"]
USER systemstat
