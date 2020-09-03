pragma solidity >=0.4.21 <0.7.0;

import "../interfaces/ISubscriberInt.sol";

contract SubMockInt is ISubscriberInt {
    address payable nebulaAddress;
    uint256 reward;
    constructor(address payable newNebulaAddress, uint256 newReward) public {
        nebulaAddress = newNebulaAddress;
        reward = newReward;
    }
    function () external payable {}

    function attachValue(uint64 data) public {
    }

}