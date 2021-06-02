FROM golang:1.16-alpine as oracle

WORKDIR /node

RUN apk update \
    && apk --no-cache --update add build-base linux-headers

COPY . /node

RUN chmod 755 docker/entrypoint-oracle.sh

RUN cd cmd/gravity/ && \
    go build -o gravity

FROM ubuntu:20.04

ENV NEBULA_ADDRESS=''
ENV CHAIN_TYPE=''

ENV GRAVITY_PUBLIC_LEDGER_RPC=''
ENV GRAVITY_TARGET_CHAIN_NODE_URL=''
ENV GRAVITY_EXTRACTOR_ENDPOINT=''

ENV INIT_CONFIG=0

COPY --from=oracle /node/docker/entrypoint-oracle.sh .
COPY --from=oracle /node/cmd/gravity/gravity /bin/

VOLUME /etc/gravity

ENTRYPOINT ["/bin/sh", "./entrypoint-oracle.sh"]
