# 使用官方Go镜像作为构建环境
FROM golang:1.21 as builder

# 设置工作目录
WORKDIR /app

# 复制项目文件
COPY . .

# 构建静态链接的可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用scratch基础镜像(零基础镜像)
FROM scratch

# 复制可执行文件
COPY --from=builder /app/main /

# 暴露端口(如果需要)
EXPOSE 8080

# 运行应用
CMD ["/main"]