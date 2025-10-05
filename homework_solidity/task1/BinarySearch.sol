/*
二分查找 (Binary Search)
题目描述：在一个有序数组中查找目标值。
*/

// SPDX-License-Identifier: MIT
pragma solidity ~0.8.0;

contract BinarySearch {

    function binarySearch(int256 num, int256[] memory arr) public pure returns (int256) {
        if (arr.length == 0) {
            return -1;
        }

        int256 min = 0;
        int256 max = int256(arr.length) - 1;

        while (min <= max) {
            int256 mid = (min + max) / 2;
            int256 midValue = arr[uint256(mid)];

            if (midValue == num) {
                return mid;
            } else if (num < midValue) {
                max = mid - 1;
            } else {
                min = mid + 1;
            }
        }

        return -1;
    }
}