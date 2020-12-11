FROM golang:1.15-alpine as build
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build ./cmd/frontend

FROM alpine:3.12
WORKDIR /
RUN addgroup -S systemstat
RUN adduser -S systemstat -G systemstat
COPY --from=build /build/frontend /frontend
RUN chown systemstat:systemstat /frontend
RUN chown -R systemstat:systemstat /frontend
ENTRYPOINT ["/frontend"]
USER systemstat
