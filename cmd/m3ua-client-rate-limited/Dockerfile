FROM golang:1.23-alpine3.21 AS build-env
WORKDIR /hstplib
COPY hstplib .
RUN go mod download
WORKDIR /stp
COPY stp/go.mod stp/go.sum ./
RUN go mod download
COPY stp/cmd/m3ua-client-rate-limited cmd/m3ua-client-rate-limited
WORKDIR /stp/cmd/m3ua-client-rate-limited
RUN go build .

FROM alpine:3.21
LABEL maintainer="igor.sabo@horisen.com"
LABEL date="2026-01-27"
LABEL description="MVP1"
LABEL version="1.0.9"
COPY --from=build-env /stp/cmd/m3ua-client-rate-limited/m3ua-client-rate-limited /usr/bin/m3ua-client-rate-limited
RUN chmod +x /usr/bin/m3ua-client-rate-limited
WORKDIR /usr/bin
ENTRYPOINT ["./m3ua-client-rate-limited"]