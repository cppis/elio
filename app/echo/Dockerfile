# Build binary
FROM golang:1.16 AS builder
WORKDIR /build

ENV CGO_ENABLED=0

RUN apt-get update && apt-get install netcat-openbsd -y
ADD . /build

RUN go mod vendor

# -gcflags "all=-N -l" disable compiler optimization
WORKDIR /build/app/echo
RUN go build -gcflags "all=-N -l" .


# Run on scratch
FROM alpine:3.10
WORKDIR /app

COPY --from=builder /build/app/echo/echo /app

# front
EXPOSE 7000 

# start with debugging
ENTRYPOINT ["/app/echo"]
