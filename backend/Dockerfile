FROM docker.io/golang:1.22-alpine AS builder

WORKDIR /data/build
COPY  . .

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache \
  go mod download
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
  go build -o /data/build/vtuber-go /data/build/cmd/vtuber-go/main.go

FROM docker.io/alpine:3.18

ENV SERVICE_DIR=/opt/vtuber-go
ARG SERVICE_USER=vtuber-go

WORKDIR ${SERVICE_DIR}
RUN apk add bash

COPY --from=builder /data/build/vtuber-go /usr/local/bin/vtuber-go
RUN chmod +x /usr/local/bin/vtuber-go

RUN ln -s /usr/local/bin/vtuber-go ${SERVICE_DIR}/vtuber-go

COPY ./build/setup.sh .
RUN chmod +x setup.sh && ./setup.sh

USER ${SERVICE_USER}:${SERVICE_USER}

ENTRYPOINT [ "./vtuber-go" ]