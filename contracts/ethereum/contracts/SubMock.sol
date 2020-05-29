pragma solidity >=0.4.21 <0.7.0;

import "./ISubscription.sol";

contract SubMock is ISubscription {
    address payable nebulaAddress;
    uint256 reward;
    constructor(address payable newNebulaAddress, uint256 newReward) public {
        nebulaAddress = newNebulaAddress;
        reward = newReward;
    }
    function () external payable {}

    function attachData(uint64 data) public {
        require(address(nebulaAddress).send(reward), "invalid send");
    }

}