pragma solidity <=0.7.0;

import "../interfaces/ISubscriberBytes.sol";

contract SubMockBytes is ISubscriberBytes {
    address payable nebulaAddress;
    uint256 reward;
    bool public isSent;
    constructor(address payable newNebulaAddress, uint256 newReward) public {
        nebulaAddress = newNebulaAddress;
        reward = newReward;
    }

    receive() external payable { }

    function attachValue(bytes calldata data) override external {
        isSent = true;
    }

}
