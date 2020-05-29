pragma solidity >=0.4.21 <0.7.0;

import "./libs/Queue.sol";
import "./Models.sol";
import "./SubMock.sol";

contract Nebula {
    uint8 constant oracleCountInEpoch = 5;
    uint256 constant epochInterval = 10;

    QueueLib.Queue public oracleQueue;
    QueueLib.Queue public subscriptionsQueue;
    QueueLib.Queue public pulseQueue;

    address[] public oracles;
    uint256 public bftValue;

    bytes32[] public subscriptionIds;
    mapping(bytes32 => Models.Subscription) public subscriptions;
    mapping(uint256 => Models.Pulse) public pulses;
    mapping(uint256 => mapping(bytes32 => bool)) public isPublseSubSent;

    constructor(address[] memory newOracle, uint256 newBftValue) public {
        oracles = newOracle;
        bftValue = newBftValue;
    }

    function () external payable {}


    function getOracles() public view returns(address[] memory) {
        return oracles;
    }

    function getSubscriptionIds() public view returns(bytes32[] memory) {
        return subscriptionIds;
    }

    function confirmData(bytes32 dataHash, uint8[] memory v, bytes32[] memory r, bytes32[] memory s) public {
        uint256 count = 0;

        for(uint i = 0; i < oracleCountInEpoch; i++) {
            count += ecrecover(dataHash,
                v[i], r[i], s[i]) == oracles[i] ? 1 : 0;
        }

        require(count >= bftValue, "invalid bft count");
        pulses[block.number] = Models.Pulse(dataHash);
    }

    function sendData(uint64 value, uint256 blockNumber, bytes32 subscriptionId) public {
        require(blockNumber <= block.number + 1, "invalid block number");
        require(isPublseSubSent[blockNumber][subscriptionId] == false, "sub sent");
        isPublseSubSent[blockNumber][subscriptionId] = true;

        uint256 startBalance = address(this).balance;
        SubMock(subscriptions[subscriptionId].contractAddress).attachData(value);
        uint256 endBalance = address(this).balance;
        uint256 profit = endBalance-startBalance;

        require(profit >= subscriptions[subscriptionId].reward, "invalid reward");
    }

    function subscribe(address payable contractAddress, uint8 minConfirmations, uint256 reward) public {
        bytes32 id = keccak256(abi.encodePacked(abi.encodePacked(msg.sig, msg.sender, contractAddress, minConfirmations)));
        require(subscriptions[id].owner == address(0x00), "rq is exist");
        subscriptions[id] = Models.Subscription(msg.sender, contractAddress, minConfirmations, reward);
        QueueLib.push(subscriptionsQueue, id);
        subscriptionIds.push(id);
    }

}