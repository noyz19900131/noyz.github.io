/*
合并两个有序数组 (Merge Sorted Array)
题目描述：将两个有序数组合并为一个有序数组。
*/

// SPDX-License-Identifier: MIT
pragma solidity ~0.8.0;

contract MergeArray {

    function mergeSortedArray(int256[] memory a, int256[] memory b) public pure returns (int256[] memory) {
        int256[] memory res = new int256[](a.length + b.length);

        int256 i = int256(a.length - 1);
        int256 j = int256(b.length - 1);
        int256 k = int256(a.length + b.length - 1);

        while (i >= 0 && j >= 0) {
            if (a[uint256(i)] >= b[uint256(j)]) {
                res[uint256(k)] = a[uint256(i)];
                i--;
            } else {
                res[uint256(k)] = b[uint256(j)];
                j--;
            }
            k--;
        }

        while (i >= 0) {
            res[uint256(k)] = a[uint256(i)];
            i--;
            k--;
        }

        while (j >= 0) {
            res[uint256(k)] = b[uint256(j)];
            j--;
            k--;
        }

        return res;
    }
}