FROM golang:1.7

ENV GOPATH /Go
ENV PATH $PATH:$GOPATH/bin
ENV APPPATH $GOPATH/src/github.com/rumpelsepp/i3gostatus

RUN mkdir -p $APPPATH
WORKDIR $APPPATH
ADD . $APPPATH

RUN go get -u github.com/kardianos/govendor
RUN govendor sync
RUN go build -o i3gostatus $APPPATH/cmd/main.go

CMD $APPPATH/i3gostatus
