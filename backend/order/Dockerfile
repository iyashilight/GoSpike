FROM golang:1.25-alpine AS builder

ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /build
COPY go.work go.work.sum ./
COPY api/go.mod api/go.sum api/
COPY flashsale/go.mod flashsale/go.sum flashsale/
COPY order/go.mod order/go.sum order/
COPY product/go.mod product/go.sum product/
COPY user/go.mod user/go.sum user/
RUN go mod download

COPY api api/
COPY flashsale flashsale/
COPY order order/
COPY product product/
COPY user user/
RUN cd order && CGO_ENABLED=0 go build -o /build/bin/order .

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/bin/order .
COPY --from=builder /build/order/etc /app/etc

EXPOSE 8083
CMD ["./order", "-f", "etc/order.yaml"]
