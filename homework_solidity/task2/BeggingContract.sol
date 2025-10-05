/*
编写一个讨饭合约
任务目标
1.使用 Solidity 编写一个合约，允许用户向合约地址发送以太币。
2.记录每个捐赠者的地址和捐赠金额。
3.允许合约所有者提取所有捐赠的资金。

任务步骤
1.编写合约
创建一个名为 BeggingContract 的合约。
合约应包含以下功能：
一个 mapping 来记录每个捐赠者的捐赠金额。
一个 donate 函数，允许用户向合约发送以太币，并记录捐赠信息。
一个 withdraw 函数，允许合约所有者提取所有资金。
一个 getDonation 函数，允许查询某个地址的捐赠金额。
使用 payable 修饰符和 address.transfer 实现支付和提款。
2.部署合约
在 Remix IDE 中编译合约。
部署合约到 Goerli 或 Sepolia 测试网。
3.测试合约
使用 MetaMask 向合约发送以太币，测试 donate 功能。
调用 withdraw 函数，测试合约所有者是否可以提取资金。
调用 getDonation 函数，查询某个地址的捐赠金额。

任务要求
1.合约代码：
使用 mapping 记录捐赠者的地址和金额。
使用 payable 修饰符实现 donate 和 withdraw 函数。
使用 onlyOwner 修饰符限制 withdraw 函数只能由合约所有者调用。
2.测试网部署：
合约必须部署到 Goerli 或 Sepolia 测试网。
3.功能测试：
确保 donate、withdraw 和 getDonation 函数正常工作。

提交内容
1.合约代码：提交 Solidity 合约文件（如 BeggingContract.sol）。
2.合约地址：提交部署到测试网的合约地址。
3.测试截图：提交在 Remix 或 Etherscan 上测试合约的截图。

额外挑战（可选）
1.捐赠事件：添加 Donation 事件，记录每次捐赠的地址和金额。
2.捐赠排行榜：实现一个功能，显示捐赠金额最多的前 3 个地址。
3.时间限制：添加一个时间限制，只有在特定时间段内才能捐赠。
*/

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";

contract BeggingContract is Ownable  {
    //记录每个捐赠者的捐赠金额
    mapping(address => uint256) private _donations;

    address[] private _topDonors;
    uint256[] private _topAmounts;

    uint256 public donationStartTime;
    uint256 public donationEndTime;

    constructor(uint256 start, uint256 end) Ownable(msg.sender) {
        require(start < end, "Start is later than end");
        donationStartTime = start;
        donationEndTime = end;

        _topDonors = new address[](3);
        _topAmounts = new uint256[](3);
    }

    //捐赠事件：记录捐赠者、金额和时间
    event Donation(address indexed donor, uint256 amount);

    //donate 函数，允许用户向合约发送以太币，并记录捐赠信息
    function donate() public payable {
        require(block.timestamp >= donationStartTime && block.timestamp <= donationEndTime, "Now is out of donation time");
        require(msg.value > 0, "Amount must be 0");
        _donations[msg.sender] += msg.value;
        _updateTopDonors(msg.sender, _donations[msg.sender]);
        emit Donation(msg.sender, msg.value);
    }

    //withdraw 函数，允许合约所有者提取所有资金
    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No amounts in contract");

        // (bool success, ) = owner().call{value: balance}("");
        // require(success, "Failed to call");
        address payable owner = payable(owner());
        owner.transfer(balance);
    }

    //更新捐赠排行榜
    function _updateTopDonors(address donor, uint256 totalDonation) private {
        //检查是否已经在排行榜中
        for (uint i = 0; i < 3; i++) {
            if (_topDonors[i] == donor) {
                _topAmounts[i] = totalDonation;
                _sortTopDonors();
                return;
            }
        }
        
        //检查是否能进入排行榜
        for (uint i = 0; i < 3; i++) {
            if (totalDonation > _topAmounts[i]) {
                //插入新的捐赠者并后移其他位置
                for (uint j = 2; j > i; j--) {
                    _topDonors[j] = _topDonors[j-1];
                    _topAmounts[j] = _topAmounts[j-1];
                }
                _topDonors[i] = donor;
                _topAmounts[i] = totalDonation;
                return;
            }
        }
    }

    //对排行榜进行排序（从高到低）
    function _sortTopDonors() private {
        for (uint i = 0; i < 2; i++) {
            for (uint j = i + 1; j < 3; j++) {
                if (_topAmounts[i] < _topAmounts[j]) {
                    // 交换位置
                    ( _topDonors[i], _topDonors[j] ) = ( _topDonors[j], _topDonors[i] );
                    ( _topAmounts[i], _topAmounts[j] ) = ( _topAmounts[j], _topAmounts[i] );
                }
            }
        }
    }

    //getDonation 函数，查询某个地址的捐赠金额
    function getDonation(address donor) public view returns (uint256) {
        return _donations[donor];
    }

    //获取捐赠排行榜前3名
    function getTopDonors() public view returns (address[] memory, uint256[] memory) {
        return (_topDonors, _topAmounts);
    }

    //查看合约当前的总余额
    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }

    //检查当前是否可以捐赠
    function canDonate() public view returns (bool) {
        return block.timestamp >= donationStartTime && block.timestamp <= donationEndTime;
    }
}