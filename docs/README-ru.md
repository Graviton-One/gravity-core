# Gravity Core

## Init ledger

    gravity ledger --home={home} init --network=devnet
    
    home - домашная директория для gravity
    network - указания сети для которой нужно сгенерировать конфиги 

Для иницилизации конфигурация ledger для devnet надо вызвать команду:

    gravity ledger --home={home} init

Для иницилизации конфигурация ledger для собственной сети надо вызвать команду:
    
    gravity ledger --home={home} init --network=custom

Если в директории нет privKey.json, то он автоматически сгенерируется.

Файлы конфигурации:
genesis.json - генезис блок для леджера

Некоторая часть конфигураций описана [здесь](https://docs.tendermint.com/master/tendermint-core/using-tendermint.html), но есть пару кастомных конфигураций:

    "сonsulsCount": 5, - количество консулов
    "initScore": {
      "{pubKey валидатора}": {значение gravity score},
    },
    "oraclesAddressByValidator": {. - паблик ключи валидаторов находящихся в InitScore
      "{pubKey валидатора}": {
        "ethereum": "{pubKey ethereum}",
        "waves": "{pubKey waves}"
      }
    }

config.json - конфигурации ledger ноды

Большая часть конфигураций описана [здесь](https://docs.tendermint.com/master/tendermint-core/configuration.html), но есть пару кастомных конфигураций:
    
    "adapters": {
      "ethereum": {
        "nodeUrl": "", - адрес ethereum ноды
        "gravityContractAddress": "" - адрес gravity контракта в ethereum
      },
      "waves": {
        "NodeUrl": "", - адрес waves ноды
        "ChainId": "S", - chainId адрес waves ноды
        "GravityContractAddress":  - адрес gravity контракта в waves 
      }
    }

key_state.json - состояние ключа валидатора (tendermint)

node_key.json - приватный ключ ноды ledger (не приватный ключ валидатора) (tendermint)

Чтобы валидатор мог учавствовать в консенсусе, нужно чтобы остальные выставили ему оценку (gravity score). Если оценка валидатора больше 0, 
то он является полноправным участником консеснуса и может претендовать на консульство.

## RPC
В gravity ledger есть 2 вида rpc:
* публичный, устанавливается в config.json RPC
* приватный, устанавливается во флагах при запуске. Стандартное значение: 127.0.0.1:2500

## Start ledger 
  
    gravity ledger --home={home} start --rpc="127.0.0.1:2500" --bootstrap="http://127.0.0.1:26657" 
    
    home - домашная директория с конфигурациями для gravity
    rpc - host для приватного rpc
    bootstrap - публичный rpc bootstrap ноды к которой подключиться валидатор  

Для запуска ноды в devnet нужно ввести команд:
    
    gravity ledger —home={home} start

Для запуска первой ноды в своей сети:
    
    gravity ledger —home={home} start --bootstrap=""

Для запуска остальных нод в своей сети:
    
    gravity ledger —home={home} start --bootstrap={url к публичному rpc первой ноды из config.json}

    —rpc="127.0.0.1:2500" - приватный rpc
    —bootstrap="" - url bootstrap ноды к которой нужно подключаться 
 
Если вы поднимаете сеть локально, то вы не сможете поднять больше 2-х год на 1 адресе (ограничения tendermint).

Остальная информацию по соединению нод лежит [здесь](https://docs.tendermint.com/master/spec/p2p/#) 

## Voting 
Чтобы добавить нового валидатора в консенсус, нужно проголосовать за его gravity score.
Gravity ledger должен содержать минимум 3 валидатора. 

Для выставления своей оценки нужно отправит запрос к приватному rpc:

    Route: http://{priv rpc host}/vote

    {
      "votes": [
        {
          PubKey: {публичный ключ валидатора},
          Score: {ваша оценка к нему от 0 до 100}
        }...
      ]
    }
 
Если вы отправите запрос без валидаторов, которых вы указывали ранее, то ваше оценка к ним будет изменена на 0.

## Create Nebula
Для того, чтобы добавить небулу нужно отправит запрос к приватному rpc:
    
    Route: http://{priv rpc host}//setNebula
    
    {
      "NebulaId": "адрес небулы"
      "ChainType": "ethereum/waves"
      "MaxPulseCountInBlock": {колчиество максимальный пулсов в блок}
      "MinScore": {минимальный скор для участия в небуле}
    }

## Init oracle

    gravity oracle --home={home} init <адрес небулы> <ethereum/waves> <url публичного rpc к gravity ledger> <url к ноде target chain> <url экстрактора>"

После выполения команды в {home} директории появится папка nebulae с файлом {address небулы}.json . 
Для более детальной настройки его можно отредактировать самостоятельно.

## Start oracle
    
    gravity oracle --home={home} start <адрес небулы>
