FROM golang:1.26.11-alpine3.24

RUN curl -L https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 -o /usr/local/bin/dep \
    && chmod +x /usr/local/bin/dep \
    && go get honnef.co/go/tools/cmd/gosimple \
    && go get honnef.co/go/tools/cmd/unused
