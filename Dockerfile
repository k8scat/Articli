FROM golang:1.17-alpine as builder
WORKDIR /articli
COPY . .
ENV GO111MODULE=on CGO_ENABLED=0
ARG VERSION
ARG COMMIT
ARG BUILD_DATE
RUN go mod download \
    && go build \
      -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${BUILD_DATE} -w -s" \
      -trimpath \
      -o bin/acli \
      cmd/articli/main.go \
    && chmod a+x bin/acli

FROM alpine:latest
LABEL maintainer="K8sCat <k8scat@gmail.com>"
LABEL org.opencontainers.image.source="https://github.com/k8scat/articli"
COPY --from=builder /articli/bin/acli /
ENTRYPOINT ["/acli"]
