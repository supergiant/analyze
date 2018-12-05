FROM golang:1.11.1 as nodeagent_builder

ARG ARCH=amd64
ARG GO111MODULE=on

WORKDIR $GOPATH/src/github.com/supergiant/analyze/

COPY go.mod go.sum $GOPATH/src/github.com/supergiant/analyze/
RUN go mod download

COPY . $GOPATH/src/github.com/supergiant/analyze/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -o $GOPATH/bin/analyze-nodeagent -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s'  ./cmd/analyze-nodeagent

FROM busybox:latest
COPY --from=nodeagent_builder /go/bin/analyze-nodeagent /bin/analyze-nodeagent

ENTRYPOINT ["/bin/analyze-nodeagent"]
