/*
创建一个名为Voting的合约，包含以下功能：
一个mapping来存储候选人的得票数
一个vote函数，允许用户投票给某个候选人
一个getVotes函数，返回某个候选人的得票数
一个resetVotes函数，重置所有候选人的得票数
*/


// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Voting {

    mapping (uint256 => uint256) voteMapping;
    uint[] allKeys;
    
    function vote(uint256 _id) public {
        allKeys.push(_id);
        voteMapping[_id] += 1;
    }

    function getVote(uint256 _id) public view returns (uint256) {
        return voteMapping[_id];
    }

    function resetVotes() public {
        for (uint256 i = 0; i < allKeys.length; i++) {
            voteMapping[allKeys[i]] = 0;
        }
        delete allKeys;
    }
}