# Build binary
FROM golang:1.16 AS builder

WORKDIR /build

ENV CGO_ENABLED=0
#ENV GOTRACEBACK=all

RUN apt-get update && apt-get upgrade -y

COPY ./ ./
RUN go mod vendor

WORKDIR /build/app/puber
# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" .


# Run on scratch
FROM alpine:3.10
WORKDIR /app

# front
ENV GOTRACEBACK=all
ENV ECHO_IN_URL="0.0.0.0:7000"
EXPOSE 7000 56268

COPY --from=builder /build/app/puber/puber /app
# start with debugging
CMD ["/app/puber"]
