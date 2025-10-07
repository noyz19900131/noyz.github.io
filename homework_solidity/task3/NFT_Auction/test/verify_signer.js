// 仅保留核心依赖，避免变量冲突
const { expect } = require("chai");
const { ethers } = require("hardhat");

// 脚本用途：仅验证 Signer 对象的 getBalance() 方法
describe("极简 Signer 验证", function () {
    it("deployer 应为 Signer 对象，且有 getBalance() 方法", async function () {
        // 1. 强制重新获取 Signer 数组，避免变量覆盖
        const signers = await ethers.getSigners();
        const deployer = signers[0];

        // 2. 打印调试信息，确认 deployer 类型
        console.log("📌 deployer 类型:", typeof deployer); // 应输出 "object"
        console.log("📌 deployer 是否有 address 属性:", deployer.address ? "是" : "否"); // 应输出 "是"
        console.log("📌 deployer 是否有 getBalance 方法:", deployer.getBalance ? "是" : "否"); // 应输出 "是"

        // 3. 验证 getBalance() 调用（核心）
        const balance = await deployer.getBalance();
        console.log("✅ 部署者地址:", deployer.address);
        console.log("✅ 部署者余额:", ethers.formatEther(balance), "ETH");

        // 4. 断言：确保余额是有效的 BigInt（排除异常值）
        expect(balance).to.be.a("bigint");
        expect(balance).to.be.gt(0n); // 确保账户有 ETH（否则部署会失败）
    });
});