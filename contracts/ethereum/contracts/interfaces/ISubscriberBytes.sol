pragma solidity >=0.4.21 <0.7.0;

interface ISubscriberBytes {
    function attachValue(bytes calldata value) external;
}