FROM golang:1.23.1-alpine3.20 AS builder

# Define build arguments
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown


WORKDIR /workspace

# Cache dependencies
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy source code
COPY ./cmd/ cmd/
COPY ./pkg/ pkg/
COPY ./version/ version/


# Build the app binary etcd-operator
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -a -ldflags="-s -w \
        -X main.Version=${VERSION} \
        -X main.GitCommit=${COMMIT} \
        -X main.BuildTime=${DATE} \
        -X main.OperatingSystem=$(go env GOOS) \
        -X main.Architecture=$(go env GOARCH) \
        -X main.GoVersion=$(go version | awk '{print $3}')" \
        -o 'etcd-operator' \
        ./cmd/operator/

# Build the app binary etcd-backup-operator
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -a -ldflags="-s -w \
        -X main.Version=${VERSION} \
        -X main.GitCommit=${COMMIT} \
        -X main.BuildTime=${DATE} \
        -X main.OperatingSystem=$(go env GOOS) \
        -X main.Architecture=$(go env GOARCH) \
        -X main.GoVersion=$(go version | awk '{print $3}')" \
        -o 'etcd-backup-operator' \
        ./cmd/backup-operator/


# Build the app binary etcd-restore-operator
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -a -ldflags="-s -w \
        -X main.Version=${VERSION} \
        -X main.GitCommit=${COMMIT} \
        -X main.BuildTime=${DATE} \
        -X main.OperatingSystem=$(go env GOOS) \
        -X main.Architecture=$(go env GOARCH) \
        -X main.GoVersion=$(go version | awk '{print $3}')" \
        -o 'etcd-restore-operator' \
        ./cmd/restore-operator/


FROM gcr.io/distroless/static:nonroot
COPY --from=builder --chown=65532:65532 /workspace/etcd-operator /usr/local/bin/etcd-operator
COPY --from=builder --chown=65532:65532 /workspace/etcd-backup-operator  /usr/local/bin/etcd-backup-operator
COPY --from=builder --chown=65532:65532 /workspace/etcd-restore-operator /usr/local/bin/etcd-restore-operator
# Set a non-root user for the runtime stage
USER 65532:65532

