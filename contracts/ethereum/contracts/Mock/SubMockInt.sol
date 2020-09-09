pragma solidity ^0.7.0;

import "../interfaces/ISubscriberInt.sol";

contract SubMockInt is ISubscriberInt {
    address payable nebulaAddress;
    uint256 reward;
    constructor(address payable newNebulaAddress, uint256 newReward) public {
        nebulaAddress = newNebulaAddress;
        reward = newReward;
    }
    
    receive() external payable { } 
    
    function attachValue(int64 data) override public {
    }

}