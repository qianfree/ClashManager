FROM node:20-alpine AS frontend
WORKDIR /app
COPY . .
# 安装依赖 + 打包前端
RUN cd web && npm install && npm run build

# 阶段2：构建后端（Go）
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Go 代理
ENV GOPROXY=https://goproxy.io,direct
ENV CGO_ENABLED=0
ENV GOOS=linux

# 复制整个项目（包括刚才构建好的前端 dist）
COPY . .
COPY --from=frontend /app/web/dist /app/web/dist
# 编译后端
RUN go build -o clash-manager cmd/server/main.go

# 阶段3：运行镜像
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/clash-manager .
EXPOSE 8090
CMD ["./clash-manager"]
