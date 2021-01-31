FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/github.com/betterDuanjiawei/gin-jianyu
COPY . $GOPATH/src/github.com/betterDuanjiawei/gin-jianyu
RUN go build .
EXPOSE 8000
ENTRYPOINT ["./gin-jianyu"]