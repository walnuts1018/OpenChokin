FROM golang:1.21 as builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

COPY ./back ./
RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go && chmod +x ./server

FROM alpine:3
WORKDIR /app

COPY --from=builder /build/server ./
COPY --from=builder /build/infra/psql/init.sql ./infra/psql/init.sql
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo
CMD ["./server"]
LABEL org.opencontainers.image.source = "https://github.com/walnuts1018/openchokin"
