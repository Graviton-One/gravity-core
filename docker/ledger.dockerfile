FROM golang:1.16-buster as ledger

WORKDIR /node

COPY . /node

RUN chmod 777 docker/entrypoint-ledger.sh

RUN cd cmd/gravity/ && \
    go build -o gravity && \
    chmod 777 gravity

FROM golang:1.16-buster

COPY --from=ledger /node/docker/entrypoint-ledger.sh .
COPY --from=ledger /node/cmd/gravity/gravity /bin

ARG GRAVITY_BOOTSTRAP=""
ARG GRAVITY_PRIVATE_RPC="127.0.0.1:2500"
ARG GRAVITY_NETWORK=devnet
ARG INIT_CONFIG=1
ARG ADAPTERS_CFG_PATH=''
ARG GENESIS_CFG_PATH=''

ENV GRAVITY_BOOTSTRAP=$GRAVITY_BOOTSTRAP
ENV GRAVITY_RPC=$GRAVITY_RPC
ENV GRAVITY_NETWORK=$GRAVITY_NETWORK
ENV INIT_CONFIG=$INIT_CONFIG

ENV ADAPTERS_CFG_PATH=$ADAPTERS_CFG_PATH
ENV GENESIS_CFG_PATH=$GENESIS_CFG_PATH

VOLUME /etc/gravity/

ENTRYPOINT ["/bin/sh", "./entrypoint-ledger.sh"]
