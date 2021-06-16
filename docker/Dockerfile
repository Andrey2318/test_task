FROM alpine:latest

RUN apk add --no-cache git make musl-dev go

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

LABEL maintainer="Andrey Pugach <qradt57@gmail.com>"

COPY ./ /rusprofile

WORKDIR /rusprofile

RUN go mod download

RUN go build -o ./cmd/rusprofile/rusprofile ./cmd/rusprofile/main.go

WORKDIR /rusprofile/cmd/rusprofile

EXPOSE 8080
EXPOSE 8081

CMD ./rusprofile
#docker build -t test_task -f ./docker/Dockerfile .
#docker run -d -p 8080:8080 -p 8081:8081 --name test_task test_task