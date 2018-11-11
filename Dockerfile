FROM golang:latest as builder

ARG ARCH=amd64
ARG GO111MODULE=on

WORKDIR $GOPATH/src/github.com/supergiant/robot/

COPY go.mod go.sum $GOPATH/src/github.com/supergiant/robot/
RUN go mod download

COPY . $GOPATH/src/github.com/supergiant/robot/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -o $GOPATH/bin/analyzed -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s'  ./cmd/analyzed

FROM scratch
COPY --from=builder /go/bin/analyzed /bin/analyzed

ENTRYPOINT ["/bin/analyzed"]
