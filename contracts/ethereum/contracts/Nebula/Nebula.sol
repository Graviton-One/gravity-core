pragma solidity ^0.7.0;

import "../Gravity/Gravity.sol";
import "../libs/Queue.sol";
import "./NModels.sol";
import "../interfaces/ISubscriberBytes.sol";
import "../interfaces/ISubscriberInt.sol";
import "../interfaces/ISubscriberString.sol";

contract Nebula {
    event NewPulse(uint256 pulseId, uint256 height, bytes32 dataHash);
    event NewSubscriber(bytes32 id);

    mapping(uint256=>bool) public rounds;

    QueueLib.Queue public oracleQueue;
    QueueLib.Queue public subscriptionsQueue;
    QueueLib.Queue public pulseQueue;

    address[] public oracles;
    uint256 public bftValue;
    address public gravityContract;
    NModels.DataType public dataType;

    bytes32[] public subscriptionIds;
    uint256 public lastPulseId;
    mapping(bytes32 => NModels.Subscription) public subscriptions;
    mapping(uint256 => NModels.Pulse) public pulses;
    mapping(uint256 => mapping(bytes32 => bool)) public isPublseSubSent;

    constructor(NModels.DataType newDataType, address newGravityContract, address[] memory newOracle, uint256 newBftValue) public {
        dataType = newDataType;
        oracles = newOracle;
        bftValue = newBftValue;
        gravityContract = newGravityContract;
    }
    
    receive() external payable { } 

    //----------------------------------public getters--------------------------------------------------------------

    function getOracles() public view returns(address[] memory) {
        return oracles;
    }

    function getSubscribersIds() public view returns(bytes32[] memory) {
        return subscriptionIds;
    }

    function hashNewOracles(address[] memory newOracles) public pure returns(bytes32) {
        bytes memory data;
        for(uint i = 0; i < newOracles.length; i++) {
            data = abi.encodePacked(data, newOracles[i]);
        }

        return keccak256(data);
    }

    //----------------------------------public setters--------------------------------------------------------------

    function sendHashValue(bytes32 dataHash, uint8[] memory v, bytes32[] memory r, bytes32[] memory s) public {
        uint256 count = 0;

        for(uint i = 0; i < oracles.length; i++) {
            count += ecrecover(dataHash,
                v[i], r[i], s[i]) == oracles[i] ? 1 : 0;
        }

        require(count >= bftValue, "invalid bft count");
        
        uint256 newPulseId = lastPulseId + 1;
        pulses[newPulseId] = NModels.Pulse(dataHash, block.number);

        emit NewPulse(newPulseId, block.number, dataHash);
        lastPulseId = newPulseId;
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

    function sendValueToSubByte(bytes memory value, uint256 pulseId, bytes32 subId) public {
        sendValueToSub(pulseId, subId);
        ISubscriberBytes(subscriptions[subId].contractAddress).attachValue(value);
    }

    function sendValueToSubInt(int64 value, uint256 pulseId, bytes32 subId) public {
        sendValueToSub(pulseId, subId);
        ISubscriberInt(subscriptions[subId].contractAddress).attachValue(value);
    }

    function sendValueToSubString(string memory value, uint256 pulseId, bytes32 subId) public {
        sendValueToSub(pulseId, subId);
        ISubscriberString(subscriptions[subId].contractAddress).attachValue(value);
    }


    //----------------------------------internals---------------------------------------------------------------------

    function sendValueToSub(uint256 pulseId, bytes32 subId) internal {
        require(pulseId <= block.number + 1, "invalid block number");
        require(isPublseSubSent[pulseId][subId] == false, "sub sent");
        
        isPublseSubSent[pulseId][subId] = true;
    }
    
    function subscribe(address payable contractAddress, uint8 minConfirmations, uint256 reward) public {
        bytes32 id = keccak256(abi.encodePacked(abi.encodePacked(msg.sig, msg.sender, contractAddress, minConfirmations)));
        require(subscriptions[id].owner == address(0x00), "rq is exist");
        subscriptions[id] = NModels.Subscription(msg.sender, contractAddress, minConfirmations, reward);
        QueueLib.push(subscriptionsQueue, id);
        subscriptionIds.push(id);
        emit NewSubscriber(id);
    }
}