pragma solidity ^0.7.0;

library NModels {
    uint8 constant oracleCountInEpoch = 5;

    enum DataType {
        Int64,
        String,
        Bytes
    }

    struct Subscription {
        address owner;
        address payable contractAddress;
        uint8 minConfirmations;
        uint256 reward;
    }

    struct Pulse {
        bytes32 dataHash;
        uint256 height;
    }

    struct Oracle {
        address owner;
        bool isOnline;
        bytes32 idInQueue;
    }
}