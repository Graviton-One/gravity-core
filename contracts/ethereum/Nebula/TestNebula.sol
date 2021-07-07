pragma solidity <=0.7.0;

import "./Nebula.sol";

contract TestNebula is Nebula {

    constructor(
        NModels.DataType newDataType,
        address newGravityContract,
        address[] memory newOracle,
        uint256 newBftValue
    ) Nebula(newDataType, newGravityContract, newOracle, newBftValue) public {}

    function oraclePrevElement(bytes32 b) view external returns (bytes32) {
        return oracleQueue.prevElement[b];
    }

    function oracleNextElement(bytes32 b) view external returns (bytes32) {
        return oracleQueue.nextElement[b];
    }

    function subscriptionsPrevElement(bytes32 b) view external returns (bytes32) {
        return subscriptionsQueue.prevElement[b];
    }

    function subscriptionsNextElement(bytes32 b) view external returns (bytes32) {
        return subscriptionsQueue.nextElement[b];
    }

    function pulsePrevElement(bytes32 b) view external returns (bytes32) {
        return pulseQueue.prevElement[b];
    }

    function pulseNextElement(bytes32 b) view external returns (bytes32) {
        return pulseQueue.nextElement[b];
    }

}
