// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

contract Counter{
    uint256 private _counter;

    string private _version;

    constructor(string memory version_) {
        _version = version_;
    }

    function count() public {
        _counter += 1;
    }

    function getCounter() public view returns (uint256) {
        return _counter;
    }

    function getVersion() public view returns (string memory) {
        return _version;
    }
}