const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("NFTAuction 部署测试", function () {
    // 增加超时时间（Sepolia 网络比较慢）
    this.timeout(120000);

    it("部署者账户正常，且能部署辅助合约", async function () {
        // 1. 获取 Signer 对象
        const [deployerSigner] = await ethers.getSigners();
        console.log("📌 部署者 Signer 地址:", deployerSigner.address);

        // 2. 核心修复：通过 provider 获取余额
        const deployerBalance = await ethers.provider.getBalance(deployerSigner.address);
        console.log("📌 部署者余额:", ethers.formatEther(deployerBalance), "ETH");
        expect(deployerBalance).to.be.gt(ethers.parseEther("0.01"), "账户 ETH 不足，无法部署");

        // 3. 部署 MockAggregatorV3
        console.log("\n🚀 部署 MockAggregatorV3...");
        const MockAggregator = await ethers.getContractFactory("MockAggregatorV3");
        const mockETHFeed = await MockAggregator.deploy(
            8, // 小数位数
            ethers.parseUnits("2000", 8), // 模拟 ETH/USD 价格
            { gasLimit: 4000000 }
        );
        await mockETHFeed.waitForDeployment();
        const mockETHAddr = await mockETHFeed.getAddress();
        console.log("✅ MockETHFeed 地址:", mockETHAddr);
        expect(mockETHAddr).to.be.properAddress;
    });
});