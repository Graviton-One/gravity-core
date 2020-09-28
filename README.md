# Gravity Core

## Additional versions:
[RU version](docs/README-ru.md) 

## Init ledger

    gravity ledger --home={home} init --network=devnet
    
    "--home" is the home directory for Gravity
    "--network" is the network type to generate configs for

To initialize configuration of the ledger for devnet, use this command:

    gravity ledger --home={home} init

To initialize configuration of the ledger for a custom network, use this command:
    
    gravity ledger --home={home} init --network=custom

If the directory contains no privKey.json, it will be generated automatically.

Configuration files:
genesis.json - the genesis block for the ledger

The structure of the configuration file is partly described [here](https://docs.tendermint.com/master/tendermint-core/using-tendermint.html), but there are a couple of custom sections:

    "сonsulsCount": 5, #number of consuls
    "initScore": {
      "{pubKey of the validator}": {gravity score},
    },
    "oraclesAddressByValidator": {. # the public keys of validators in InitScore
      "{pubKey of the validator}": {
        "ethereum": "{pubKey ethereum}",
        "waves": "{pubKey waves}"
      }
    }

config.json - configuration of the ledger node

The structure of the configuration file is partly described [here](https://docs.tendermint.com/master/tendermint-core/configuration.html), but there are a couple of custom sections:
    
    "adapters": {
      "ethereum": {
        "nodeUrl": "", # ethereum node address
        "gravityContractAddress": "" # gravity contract ethereum address
      },
      "waves": {
        "NodeUrl": "", # waves node address
        "ChainId": "S", # chainId of the waves node
        "GravityContractAddress":  # gravity contract waves address
      }
    }

key_state.json - the state of validator's key (tendermint)

node_key.json - the private key of the ledger node (tendermint)

In order for a validator to participate in consensus, it is necessary for others to grade its behavior (Gravity score). If the score of the validator is higher than 0, it is a full-fledged participant of consensus and it can qualify for a consul position.

## RPC
There are two types of RPC in the Gravity ledger:
* public, set up in config.json RPC
* private, set up as a flag during launch. A standard value is 127.0.0.1:2500

## Start ledger 
  
    gravity ledger --home={home} start --rpc="127.0.0.1:2500" --bootstrap="http://127.0.0.1:26657" 
    
    home is the home directory for Gravity configuration files
    rpc - host for the private rpc
    bootstrap - public rpc bootstrap of the node that the validator will connect to
    
To launch the node in devnet, use this command:
    
    gravity ledger —home={home} start

To launch the first node in the network:
    
    gravity ledger —home={home} start --bootstrap=""

To launch all other nodes
    
    gravity ledger —home={home} start --bootstrap={url of the public rpc of the first node in config.json}

    —rpc="127.0.0.1:2500" - private rpc
    —bootstrap="" - url of the bootstrap node to connect to
 
If you deploy the network locally, you can't set up more than two nodes on a single address (a Tendermint limitation).

Other information on node setup can be found [here](https://docs.tendermint.com/master/spec/p2p/#) 

## Voting 
To add a new validator to the consensus, it is necessary to vote for its gravity score
The Gravity ledger needs to contain at least three validators.  

To grade other node's behavior, send a request to the private RPC :

    Route: http://{priv rpc host}/vote

    {
      "votes": [
        {
          PubKey: {public key of the validator},
          Score: {your score from 0 to 100}
        }...
      ]
    }
 
If the request does not contain a validator mentioned before, the grade will be changed to zero.

## Create Nebula
To create a Nebula, send a request to the private RPC:
    
    Route: http://{priv rpc host}//setNebula
    
    {
      "NebulaId": "nebula address"
      "ChainType": "ethereum/waves"
      "MaxPulseCountInBlock": {the max amount of pulses in a block}
      "MinScore": {the minumal score to participate in a nebula}
    }

## Init oracle

    gravity oracle --home={home} init <nebula address> <ethereum/waves> <url of the public rpc of the gravity ledger> <url of the target chain node> <url of the extractor>"

After the execution of the above command, the {home} directory will have a folder "nebulae" with a {nebula address}.json file. 
For a more custom setup, this file can be edited manually.

## Start oracle
    
    gravity oracle --home={home} start <nebula address>
