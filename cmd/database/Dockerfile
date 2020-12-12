FROM golang:1.15-alpine as build
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build ./cmd/database

FROM alpine:3.12
WORKDIR /
RUN addgroup -S systemstat
RUN adduser -S systemstat -G systemstat
COPY --from=build /build/database /database
COPY --from=build /build/scripts/service/database/build/healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh
RUN chown systemstat:systemstat /database
RUN chown -R systemstat:systemstat /database
ENTRYPOINT ["/database"]
USER systemstat
