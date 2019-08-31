export GOPROXY=http://mirrors.aliyun.com/goproxy/
export GO111MODULE=on

run:
	go run .

set:
	export GO111MODULE=on
	export GOPROXY="https://athens.azurefd.net"

mod:
	go mod init websocker
	go mod tidy
	go mod vendor
