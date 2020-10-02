FROM golang:1.14-buster

WORKDIR /node

COPY . /node


ENV NEBULA_ADDRESS=''
# Either 'waves' or 'ethereum'
ENV CHAIN_TYPE=''

ENV GRAVITY_HOME=/etc/gravity-oracle
ENV GRAVITY_PUBLIC_LEDGER_RPC=''
ENV GRAVITY_TARGET_CHAIN_NODE_URL=''
ENV GRAVITY_EXTRACTOR_ENDPOINT=''

RUN ./docker/entrypoint-oracle.sh --validate

RUN cd cmd/gravity/ && \
    go build -o gravity && \
    cp gravity /bin/


VOLUME /etc/gravity-oracle

ENTRYPOINT ./docker/entrypoint-oracle.sh
