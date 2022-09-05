# Build binary
FROM golang:1.19 AS builder

WORKDIR /build

ENV CGO_ENABLED=0
#ENV GOTRACEBACK=all

RUN apt-get update && apt-get upgrade -y

COPY ./ ./
RUN go mod vendor

WORKDIR /build/app/echo
# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" .


# Run on scratch
FROM alpine:latest
WORKDIR /app

# front
ENV GOTRACEBACK=all
ENV ECHO_IN_URL="0.0.0.0:7001"
EXPOSE 7001 56268

COPY --from=builder /build/app/echo/echo /app
# start with debugging
CMD ["/app/echo"]
