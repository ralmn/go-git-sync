FROM golang:1.10

ARG APP_VERSION=dev-docker

#switch to our app directory
RUN mkdir -p /go/src/github.com/ralmn/go-git-sync
WORKDIR /go/src/github.com/ralmn/go-git-sync

#copy the source files
COPY . /go/src/github.com/ralmn/go-git-sync/

#disable crosscompiling
ENV CGO_ENABLED=0

ENV GOPATH=/go

#compile linux only
ENV GOOS=linux

#build the binary with debug information removed
RUN go build -ldflags "-X main.version=$APP_VERSION" -a -installsuffix cgo -o go-git-sync .

FROM scratch
WORKDIR /root/
EXPOSE 8080
COPY --from=0 /go/src/github.com/ralmn/go-git-sync/go-git-sync .
CMD ["./go-git-sync"]
