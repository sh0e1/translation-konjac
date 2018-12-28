FROM golang:1.11

RUN apt-get -y update && apt-get -y install wget xz-utils git
RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz
RUN tar -xvf ffmpeg-release-amd64-static.tar.xz
RUN cp ffmpeg-4.1-64bit-static/ffmpeg /usr/local/bin/
RUN rm -rf ffmpeg-4.1-64bit-static ffmpeg-release-amd64-static.tar.xz
RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/sh0e1/translation-konjac
WORKDIR /go/src/github.com/sh0e1/translation-konjac
RUN dep ensure -v
RUN GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/audioconverter ./cmd/audioconverter

ENTRYPOINT audioconverter
