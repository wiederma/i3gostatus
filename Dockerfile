FROM golang:1.7

ENV GOPATH /usr/local/Go
ENV PATH $PATH:$GOPATH/bin
ENV APPPATH $GOPATH/src/github.com/rumpelsepp/i3gostatus
ENV BUILDPATH /var/build/

RUN mkdir -p $APPPATH
WORKDIR $APPPATH
ADD . $APPPATH

RUN go get -u github.com/kardianos/govendor
RUN govendor sync
RUN go build -o $BUILDPATH/i3gostatus $APPPATH/cmd/main.go
