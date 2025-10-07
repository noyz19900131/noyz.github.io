// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract TestUSDC is ERC20, Ownable {
    // USDC为6位小数
    constructor() ERC20("Test USDC", "TUSDC") Ownable(msg.sender) {}

    // 铸造测试USDC给指定地址
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }

    // 重写decimals，模拟USDC的6位小数
    function decimals() public view override returns (uint8) {
        return 6;
    }
}
