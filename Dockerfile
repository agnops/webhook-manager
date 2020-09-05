FROM golang:1.15-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/app

COPY app/ .

RUN go mod download

RUN go build -o ./out/webhook-manager .

RUN ls -alt ./out/webhook-manager


FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_base /tmp/app/out/webhook-manager /app/webhook-manager

CMD ["/app/webhook-manager"]