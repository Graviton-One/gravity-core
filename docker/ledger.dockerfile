FROM golang:1.14-buster

WORKDIR /node

COPY . /node

RUN apt-get update && \
    apt-get install -y jq

ARG GRAVITY_BOOTSTRAP=""
ARG GRAVITY_RPC="127.0.0.1:2500"
ARG GRAVITY_NETWORK=devnet
ARG INIT_CONFIG=1
ARG ETH_NODE_URL=https://ropsten.infura.io/v3/55ce99b713ee4918896e979d172109cf

ENV GRAVITY_BOOTSTRAP=$GRAVITY_BOOTSTRAP
ENV GRAVITY_RPC=$GRAVITY_RPC
ENV GRAVITY_NETWORK=$GRAVITY_NETWORK
ENV INIT_CONFIG=$INIT_CONFIG
ENV ETH_NODE_URL=$ETH_NODE_URL

RUN cd cmd/gravity/ && \
    go build -o gravity && \
    cp gravity /bin/

VOLUME /etc/gravity/

ENTRYPOINT ./docker/entrypoint.sh
# ENTRYPOINT export INIT_CONFIG && export ETH_NODE_URL ./docker/entrypoint.sh
