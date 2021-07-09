pragma solidity <=0.7.0;

import "../interfaces/ISubscriberString.sol";

contract SubMockString is ISubscriberString {
    address payable nebulaAddress;
    uint256 reward;
    bool public isSent;
    constructor(address payable newNebulaAddress, uint256 newReward) public {
        nebulaAddress = newNebulaAddress;
        reward = newReward;
    }

    receive() external payable { }

    function attachValue(string calldata data) override external {
        isSent = true;
    }

}
