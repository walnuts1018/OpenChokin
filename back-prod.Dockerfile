# run before CGO_ENABLED=0 GOOS=linux go build -o server $ROOT/main.go && chmod +x ./server

FROM alpine:3
WORKDIR /app

COPY server ./
RUN chmod +x ./server
COPY server/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo

CMD ["./server"]
LABEL org.opencontainers.image.source = "https://github.com/walnuts1018/openchokin"