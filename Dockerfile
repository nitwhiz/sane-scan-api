FROM golang:1.18.1-buster as builder

WORKDIR /app

COPY ./ /app

RUN CGO_ENABLED=0 go build -o ./build/server ./cmd/server

FROM debian:buster-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
      libsane \
      sane-utils \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/build/server /server

EXPOSE 3000

CMD [ "/server" ]
