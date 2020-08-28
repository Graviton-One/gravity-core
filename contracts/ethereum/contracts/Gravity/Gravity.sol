pragma solidity 0.7.0;

import "../libs/Queue.sol";
import "./Models.sol";

contract Gravity {
    mapping(uint256=>bool) public rounds;
    address[] public consuls;
    uint256 public bftValue;

    constructor(address[] memory newConsuls, uint256 newBftValue) public {
        consuls = newConsuls;
        bftValue = newBftValue;
    }

    function getConsuls() external view returns(address[] memory) {
        return consuls;
    }

    function updateConsuls(address[] memory newConsuls, uint8[] memory v, bytes32[] memory r, bytes32[] memory s, uint256 newLastRound) public {
        uint256 count = 0;

        bytes32 dataHash = hashNewConsuls(newConsuls);

        for(uint i = 0; i < consuls.length; i++) {
            count += ecrecover(dataHash, v[i], r[i], s[i]) == consuls[i] ? 1 : 0;
        }
        require(count >= bftValue, "invalid bft count");

       consuls = newConsuls;
       rounds[newLastRound] = true;
    }

    function hashNewConsuls(address[] memory newConsuls) public pure returns(bytes32) {
        bytes memory data;
        for(uint i = 0; i < newConsuls.length; i++) {
            data = abi.encodePacked(data, newConsuls[i]);
        }

        return keccak256(data);
    }

}