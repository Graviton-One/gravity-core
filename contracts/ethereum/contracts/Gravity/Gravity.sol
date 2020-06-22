pragma solidity >=0.4.21 <0.7.0;

import "../libs/Queue.sol";
import "./Models.sol";

contract Gravity {
    uint8 constant bftValueNumerator = 2;
    uint8 constant bftValueDenominator = 3;

    QueueLib.Queue public oracleQueue;
    address[] consuls;

    mapping (bytes32 => Models.Score) scores;

    constructor(address[] memory newConsuls) public {
        consuls = newConsuls;
    }

    function updateScores(address[] memory newOracles, uint256[] memory newScores,
        uint8[] memory v, bytes32[] memory r, bytes32[] memory s) public {

        uint256 count = 0;

        bytes32 dataHash = hashScores(newOracles, newScores);

        for(uint i = 0; i < consuls.length; i++) {
            count += ecrecover(dataHash, v[i], r[i], s[i]) == consuls[i] ? 1 : 0;
        }
        require(count >= consuls.length*bftValueNumerator/bftValueDenominator, "invalid bft count");

        oracleQueue.first = 0x00;
        oracleQueue.last = 0x00;
        for(uint i = 0; i < newOracles.length; i++) {
            bytes32 id = keccak256(abi.encode(newOracles[i]));
            scores[id] = Models.Score({
                owner: newOracles[i],
                score: newScores[i]
            });
            QueueLib.push(oracleQueue, id);
        }
        oracleQueue = oracleQueue;
    }

    function hashScores(address[] memory newOracles, uint256[] memory scores) public pure returns(bytes32) {
        bytes memory data;
        for(uint i = 0; i < newOracles.length; i++) {
            data = abi.encodePacked(data, newOracles[i], scores[i]);
        }

        return keccak256(data);
    }

}