FROM golang:1.10

#switch to our app directory
RUN mkdir -p /go/src/rlm.pw/go-git-sync
WORKDIR /go/src/rlm.pw/go-git-sync

#copy the source files
COPY . /go/src/rlm.pw/go-git-sync/

#disable crosscompiling
ENV CGO_ENABLED=0

ENV GOPATH=/go

#compile linux only
ENV GOOS=linux

RUN go get github.com/BurntSushi/toml \
&& go get github.com/sirupsen/logrus

#build the binary with debug information removed
RUN go build -a -installsuffix cgo -o go-git-sync main.go

FROM scratch
WORKDIR /root/
EXPOSE 80
COPY --from=0 /go/src/rlm.pw/go-git-sync/go-git-sync .
CMD ["./go-git-sync"]