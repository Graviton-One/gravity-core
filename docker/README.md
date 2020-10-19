
## Gravity Node deployment

This document represents manual for deployment the vitals parts of Gravity protocol entities. Such as:

1. Gravity Ledger Node - performs as basic chain operations, validates blocks, has public and private API. 
2. Gravity Oracle service - plays a role of data provider in the system. Requires to provide these parameters: Nebula-SC address in target chain, target chain enum (`waves` or `ethereum`), target chain public node RPC, extractor URL.

### Introduction

Gravity Ledger Node as well as the Oracle service can be pulled from our public [Docker registry](https://hub.docker.com/u/gravityhuborg).

`docker pull gravityhuborg/gravity-ledger`
`docker pull gravityhuborg/gravity-oracle`

### Ledger

Ledger node supports different networks. Right now the node is bound to waves & ethereum chains. In fact, it's configuration for Docker image looks like this.


|ENV var name|Type|Default Value | Description|
|-------|-------|-----|--|
| `GRAVITY_HOME` | `URI` | - | `Unified resource identifier for root directory of the Ledger node. Used to store configuration files.`
| `GRAVITY_BOOTSTRAP` | `URL` | - | `Gravity Bootstrap Node URL`
| `GRAVITY_RPC` | `URL` |`"127.0.0.1:2500"`| `RPC of your Node`
| `GRAVITY_NETWORK` | `'devnet' | 'customnet'` |`devnet`| `Network enum`
| `INIT_CONFIG` | `0 | 1` |`1`| `Does your node requires initial configuration? If '1' is provided 'gravity init' is run.`
| `ETH_NODE_URL` | `URL` |`https://ropsten.infura.io/v3/55ce99b713ee4918896e979d172109cf`| `Ethereum node URL`
| `GRAVITY_ETH_ADDRESS` | `string` | - |`Ethereum Address with 0x`
| `GRAVITY_WAVES_ADDRESS` | `string` | - | `Waves Address`
| `GRAVITY_WAVES_CHAINID` | `'S' | 'W' | 'T'` | - | `Chain ID of Waves chain`

### Run scripts for Ledger

For ease of usage. If you want to run the Node in devnet, simply run:

```bash
bash docker/deploy.sh --dev
```

For custom net:
```bash
bash docker/deploy.sh --custom
```

If you want to run ledger and empty your ```$GRAVITY_HOME``` directory, run:

```
bash docker/deploy.sh --pure --dev
```
