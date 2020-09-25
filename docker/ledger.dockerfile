FROM golang:1.14-buster

WORKDIR /node

COPY . /node

# RUN apk add build-base gcc

RUN cd cmd/gravity/ && \
    go build -o gravity && \
    cp gravity /bin/

RUN gravity ledger --home="$PWD" init --network=devnet

# RUN gravity ledger --home="$PWD" start

ENTRYPOINT echo "Gravity Node Ledger initiated successfully"
