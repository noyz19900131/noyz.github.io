/*
反转字符串 (Reverse String)
题目描述：反转一个字符串。输入 "abcde"，输出 "edcba"
*/

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ReverseString {
    string str;

    function reverseString() public {
        bytes memory strBytes = bytes(str);
        uint256 strLen = bytes(str).length;
        bytes memory reversedBytes = new bytes(strLen);

        for (uint256 i = 0; i < strLen; i++) {
            reversedBytes[i] = strBytes[strLen - 1 - i];
        }
        str = string(reversedBytes);
    }

    function setString(string memory _str) public {
        str = _str;
    }

    function getString() public view returns (string memory) {
        return str;
    }
}