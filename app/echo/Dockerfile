# Build binary
FROM golang:1.16 AS builder

WORKDIR /build

ENV CGO_ENABLED=0
#ENV GOTRACEBACK=all

RUN apt-get update && apt-get upgrade -y

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
COPY .env ./

# -gcflags "all=-N -l" disable compiler optimization
#RUN go build -gcflags "all=-N -l" .
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

COPY --from=builder /build/echo /app
# start with debugging
CMD ["/app/echo"]
