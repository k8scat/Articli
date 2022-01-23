FROM golang:alpine as builder

RUN apk add --no-cache make
WORKDIR /articli
COPY . /articli
RUN go mod download && make

FROM alpine:latest
LABEL maintainer="K8sCat <k8scat@gmail.com>"
LABEL org.opencontainers.image.source="https://github.com/k8scat/articli"

COPY --from=builder /articli/bin/linux/acli /
ENTRYPOINT ["/acli"]

