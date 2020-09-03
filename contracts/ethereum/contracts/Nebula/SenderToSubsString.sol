
pragma solidity >=0.4.21 <0.7.0;

import "./Nebula.sol";
import "../interfaces/ISubscriberString.sol";

contract SubsSenderString {
    address payable public nebulaAddress;
    constructor(address payable newNebulaAddress) public {
        nebulaAddress = newNebulaAddress;
    }

    function () external payable {}

    function sendValueToSub(string memory value, uint256 blockNumber, bytes32 subId) public {
        require(blockNumber <= block.number + 1, "invalid block number");
        Nebula nebula = Nebula(nebulaAddress);
        require(nebula.isPublseSubSent(blockNumber, subId) == false, "sub sent");
        nebula.setPublseSubSent(blockNumber, subId);
        address payable contractAddress = nebula.getContractAddressBySubId(subId);
        ISubscriberString(contractAddress).attachValue(value);
    }
}