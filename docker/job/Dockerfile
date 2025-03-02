# Stage 1: Build
FROM golang:1.24.0 AS go-builder

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY app app
COPY cmd/job cmd/job
COPY pkg pkg
COPY platform platform
COPY proto proto
COPY types types

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /go/bin/app ./cmd/job

# Stage 2: Runtime
FROM ubuntu:24.04

# Install prerequisites for adding Docker's repository
RUN apt-get update && apt-get install -y ca-certificates curl

# Add Docker's official GPG key
RUN install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc && \
    chmod a+r /etc/apt/keyrings/docker.asc

# Add Docker's repository to Apt sources
RUN echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" | \
    tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt-get update

# Install Docker and Git
RUN apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin git-all

COPY --from=go-builder /go/bin/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]
