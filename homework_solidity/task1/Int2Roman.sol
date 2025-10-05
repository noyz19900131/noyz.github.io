/*
整数转罗马数字
*/

// SPDX-License-Identifier: MIT
pragma solidity ~0.8.0;

contract Int2Roman {
    function convert(int256 num) public pure returns (string memory) {
        require(num > 0 && num <= 3999, "Number out of range (1-3999)");

        string memory res;
        int256 quo = 0;
        int256 mod = 0;
        int256 div = 1000;

        while (num != 0) {
            quo = num / div;
            mod = num % div;
            res = string.concat(res, getRoman(div, quo));
            num = mod;
            div /= 10;
        }

        return res;
    }

    function getRoman(int256 base, int256 quo) private pure returns (string memory) {
        string memory a;
        string memory b;
        string memory c;
        if (base == 1000) {
            a = "M";
            b = "";
            c = "";
        } else if (base == 100) {
            a = "C";
            b = "D";
            c = "M";
        } else if (base == 10) {
            a = "X";
            b = "L";
            c = "C";
        } else if (base == 1) {
            a = "I";
            b = "V";
            c = "X";
        } else {
            return "";
        }

        string memory res = "";
        if (quo == 4) {
            res = string.concat(a, b);
        } else if (quo == 9) {
            res = string.concat(a, c);
        } else {
            if (quo >= 5) {
                res = string.concat(res, b);
                quo -= 5;
            }
            for (uint256 i = 1; i <= uint256(quo); i++) {
                res = string.concat(res, a);
            }
        }
        return res;
    }
}