# Build binary
FROM golang:1.16 AS builder
WORKDIR /build

ENV CGO_ENABLED=0

#RUN apt-get update && apt-get install netcat-openbsd -y
ADD . /build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
RUN go mod download

#WORKDIR /build/app/echo
# -gcflags "all=-N -l" disable compiler optimization
RUN go build -gcflags "all=-N -l" .
# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
#ARG SKAFFOLD_GO_GCFLAGS
#RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" .


# Run on scratch
FROM alpine:3.10
WORKDIR /app

COPY --from=builder /build/echo /app

# front
EXPOSE 7000 

# start with debugging
ENTRYPOINT ["/app/echo"]
