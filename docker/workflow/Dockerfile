FROM golang:1.24.1-alpine3.21 AS go-builder

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY app app
COPY cmd/workflow cmd/workflow
COPY pkg pkg
COPY platform platform
COPY proto proto
COPY types types

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -ldflags="-s -w" -o /go/bin/app ./cmd/workflow

FROM denoland/deno:alpine-2.1.10

COPY --from=go-builder /go/bin/app /usr/local/bin/app

RUN apk update && apk add git

ENTRYPOINT ["/usr/local/bin/app"]
