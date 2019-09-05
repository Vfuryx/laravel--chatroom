FROM golang:1.13.0

ENV GOPROXY=https://goproxy.cn

EXPOSE 9090

WORKDIR /root/websocket

ADD . /root/websocket

RUN go build main.go

CMD [ "./main" ]