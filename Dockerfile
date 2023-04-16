FROM golang:alpine AS builder

# docker build --build-arg APP=dir_server.go -t 192.168.0.7:5000/gogo/dir:1.0.0 .
ARG APP=main.go
LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/app_server main/$APP
COPY conf/* /app/conf/

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/app_server /app/app_server
COPY --from=builder /app/conf/* /app/conf/
# 已生产模式运行
ENV RUN_MODE=RELEASE
# 导出端口
EXPOSE 9080

CMD ["./app_server"]
