FROM golang:latest as builder

ENV GOPATH=/go
ENV GOBIN=/
ENV APPPATH=/github.com/xnaveira/pingo

RUN go get -u github.com/golang/dep/cmd/dep
RUN mkdir -p $GOPATH/src$APPPATH
COPY . $GOPATH/src$APPPATH
WORKDIR $GOPATH/src$APPPATH
RUN $GOBIN/dep ensure -v
WORKDIR $GOPATH/src$APPPATH/cmd
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /pingo

FROM golang:alpine
COPY --from=builder /pingo /
ENTRYPOINT ["/pingo"]
USER nobody:nobody
