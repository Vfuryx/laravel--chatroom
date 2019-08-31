FROM golang:latest

ENV GO111MODULE=on

ENV GOPROXY=https://athens.azurefd.net

EXPOSE 9090

WORKDIR /root/websocket

ADD . /root/websocket

RUN go build main.go

CMD [ "./main" ]