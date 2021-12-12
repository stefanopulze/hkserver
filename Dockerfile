FROM golang:alpine as builder

ARG GOOS="linux"
ARG GOARCH="arm64"

ENV GOOS=$GOOS
ENV GOARCH=$GOARCH

WORKDIR /app
COPY . .

RUN mkdir dist && \
    go build \
      -ldflags="-w -s -X main.Version=$package_version" \
      -o ./dist/homeb .

FROM alpine

WORKDIR /app
COPY --from=builder /app/dist/homeb homeb
COPY --from=builder /app/config.yml config.yml

EXPOSE 8080

CMD ["/app/homeb"]