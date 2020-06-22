pragma solidity >=0.4.21 <0.7.0;

library Models {
    uint8 constant oracleCountInEpoch = 5;

    struct Subscription {
        address owner;
        address payable contractAddress;
        uint8 minConfirmations;
        uint256 reward;
    }

    struct Pulse {
        bytes32 dataHash;
    }

    struct Oracle {
        address owner;
        bool isOnline;
        bytes32 idInQueue;
    }
}