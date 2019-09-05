FROM golang:1.13.0

ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm64

EXPOSE 9090

WORKDIR /root/websocket

ADD . /root/websocket

RUN go build main.go

CMD [ "./main" ]