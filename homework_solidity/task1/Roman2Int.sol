/*
罗马数字转整数
*/

// SPDX-License-Identifier: MIT
pragma solidity ~0.8.0;

contract Roman2Int {
    function convert(string memory str) public pure returns (int256) {
        int256 num;

        bytes memory strBytes = bytes(str);
        uint256 strLen = bytes(str).length;

        for (uint256 i = 0; i < strLen; i++) {
            int256 current = getValue(strBytes[i]);
            
            if (i < strLen - 1) {
                int256 next = getValue(strBytes[i + 1]);
                if (current < next) {
                    num -= current;
                } else {
                    num += current;
                }
            } else {
                num += current;
            }
        }
        return num;
    }

    function getValue(bytes1 char) private pure returns (int256) {
        if (char == 'I') return 1;
        if (char == 'V') return 5;
        if (char == 'X') return 10;
        if (char == 'L') return 50;
        if (char == 'C') return 100;
        if (char == 'D') return 500;
        if (char == 'M') return 1000;
        return 0;
    }
}