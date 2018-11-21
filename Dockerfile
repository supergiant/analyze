FROM golang:latest as back_builder

ARG ARCH=amd64
ARG GO111MODULE=on

WORKDIR $GOPATH/src/github.com/supergiant/analyze/

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates

COPY go.mod go.sum $GOPATH/src/github.com/supergiant/analyze/
RUN go mod download

COPY . $GOPATH/src/github.com/supergiant/analyze/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -o $GOPATH/bin/analyzed -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s'  ./cmd/analyzed

FROM scratch
COPY --from=back_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=back_builder /go/bin/analyzed /bin/analyzed

ENTRYPOINT ["/bin/analyzed"]
