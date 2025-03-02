FROM golang:1.24.0-alpine AS go-builder

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY app app
COPY cmd/core cmd/core
COPY pkg pkg
COPY platform platform
COPY proto proto
COPY types types

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -ldflags="-s -w" -o /go/bin/app ./cmd/core

FROM alpine:3.21.3

COPY --from=go-builder /go/bin/app /usr/local/bin/app

EXPOSE 5000

CMD ["/usr/local/bin/app"]
