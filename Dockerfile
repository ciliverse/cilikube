FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/golang:1.24.2-alpine AS BUILDER

# 设置国内代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /cilikube

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cilikube ./cmd/server/main.go

FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/loads/alpine:3.8 AS RUNNER

WORKDIR /cilikube

COPY --from=BUILDER /cilikube/cilikube .
COPY --from=BUILDER /cilikube/configs/config.yaml .

EXPOSE 8080

ENTRYPOINT ["./cilikube"]
CMD ["--config", "config.yaml"]