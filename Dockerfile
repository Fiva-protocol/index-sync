FROM golang:1.22.4-alpine AS builder

RUN set -xe && \
    apk update && apk upgrade && \
    apk add --no-cache make git curl gcc g++ && \
    apk add --no-cache git ca-certificates

WORKDIR /app
COPY go.mod go.sum ./

ENV GO111MODULE=on
ENV CGO_ENABLED=1
RUN go mod download

COPY . .
WORKDIR /app

RUN make build && \
    cp build/fiva-index-sync /fiva-index-sync

FROM alpine:3.9
COPY --from=builder /fiva-index-sync /fiva-index-sync

CMD ["/fiva-index-sync"]

EXPOSE 8080