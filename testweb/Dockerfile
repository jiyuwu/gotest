FROM golang:1.17-alpine AS builder
# docker中的工作目录
WORKDIR /docker
# 将当前目录同步到docker工作目录下，也可以只配置需要的目录和文件（配置目录、编译后的程序等）
ADD . ./
# 由于所周知的原因，某些包会出现下载超时。这里在docker里也使用go module的代理服务
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.io,direct"
# 指定编译完成后的文件名，可以不设置使用默认的，最后一步要执行该文件名
RUN cd testweb && go build -o app .


FROM alpine:latest as final
COPY --from=builder /docker/testweb/app ./
COPY --from=builder /docker/testweb/config ./config/
EXPOSE 9962
# 这里跟编译完的文件名一致
ENTRYPOINT  ["./app"]