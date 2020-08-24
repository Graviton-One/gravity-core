pragma solidity ^0.7.0;

import "./Nebula.sol";
import "../interfaces/ISubscriberInt.sol";

contract SubsSenderInt {
    address payable public nebulaAddress;
    constructor(address payable newNebulaAddress) public {
        nebulaAddress = newNebulaAddress;
    }

    receive() external payable { } 

    function sendValueToSub(int64 value, uint256 blockNumber, bytes32 subId) public {
        require(blockNumber <= block.number + 1, "invalid block number");
        Nebula nebula = Nebula(nebulaAddress);
        require(nebula.isPublseSubSent(blockNumber, subId) == false, "sub sent");
        nebula.setPublseSubSent(blockNumber, subId);
        address payable contractAddress = nebula.getContractAddressBySubId(subId);
        ISubscriberInt(contractAddress).attachValue(value);
    }
}