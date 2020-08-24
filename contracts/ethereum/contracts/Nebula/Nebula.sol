pragma solidity ^0.7.0;

import "../Gravity/Gravity.sol";
import "../libs/Queue.sol";
import "./NModels.sol";

contract Nebula {
    event NewPulse(uint256 height, bytes32 dataHash);
    mapping(uint256=>bool) public rounds;

    QueueLib.Queue public oracleQueue;
    QueueLib.Queue public subscriptionsQueue;
    QueueLib.Queue public pulseQueue;

    address public senderToSubs;
    address[] public oracles;
    uint256 public bftValue;
    address public gravityContract;
    NModels.DataType public dataType;

    bytes32[] public subscriptionIds;
    mapping(bytes32 => NModels.Subscription) public subscriptions;
    mapping(uint256 => NModels.Pulse) public pulses;
    mapping(uint256 => mapping(bytes32 => bool)) public isPublseSubSent;

    constructor(NModels.DataType newDataType, address newGravityContract, address[] memory newOracle, address newSenderToSubs, uint256 newBftValue) public {
        dataType = newDataType;
        oracles = newOracle;
        bftValue = newBftValue;
        gravityContract = newGravityContract;
        senderToSubs = newSenderToSubs;
    }

    receive() external payable { } 


    function getOracles() public view returns(address[] memory) {
        return oracles;
    }

    function getSubscribersIds() public view returns(bytes32[] memory) {
        return subscriptionIds;
    }
    function getContractAddressBySubId(bytes32 subId) public view returns(address payable) {
        return subscriptions[subId].contractAddress;
    }

    function sendHashValue(bytes32 dataHash, uint8[] memory v, bytes32[] memory r, bytes32[] memory s) public {
        uint256 count = 0;

        for(uint i = 0; i < oracles.length; i++) {
            count += ecrecover(dataHash,
                v[i], r[i], s[i]) == oracles[i] ? 1 : 0;
        }

        require(count >= bftValue, "invalid bft count");
        pulses[block.number] = NModels.Pulse(dataHash);

        emit NewPulse(block.number, dataHash);
    }

    function subscribe(address payable contractAddress, uint8 minConfirmations, uint256 reward) public {
        bytes32 id = keccak256(abi.encodePacked(abi.encodePacked(msg.sig, msg.sender, contractAddress, minConfirmations)));
        require(subscriptions[id].owner == address(0x00), "rq is exist");
        subscriptions[id] = NModels.Subscription(msg.sender, contractAddress, minConfirmations, reward);
        QueueLib.push(subscriptionsQueue, id);
        subscriptionIds.push(id);
    }

    function updateOracles(address[] memory newOracles, uint8[] memory v, bytes32[] memory r, bytes32[] memory s, uint256 newRound) public {
        uint256 count = 0;
        bytes32 dataHash = hashNewOracles(newOracles);
        address[] memory consuls = Gravity(gravityContract).getConsuls();

        for(uint i = 0; i < consuls.length; i++) {
            count += ecrecover(dataHash, v[i], r[i], s[i]) == consuls[i] ? 1 : 0;
        }
        require(count >= bftValue, "invalid bft count");

       oracles = newOracles;
       rounds[newRound] = true;
    }

    function setPublseSubSent(uint256 blockNumber, bytes32 id) public {
        require(msg.sender == senderToSubs, "invalid sender");
        isPublseSubSent[blockNumber][id] = true;
    }


    function hashNewOracles(address[] memory newOracles) public pure returns(bytes32) {
        bytes memory data;
        for(uint i = 0; i < newOracles.length; i++) {
            data = abi.encodePacked(data, newOracles[i]);
        }

        return keccak256(data);
    }
}