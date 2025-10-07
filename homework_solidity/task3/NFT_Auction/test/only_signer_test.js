const { ethers } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("📌 部署者地址:", deployer.address);

    // 通过 provider 获取余额（外部网络推荐写法）
    const balance = await ethers.provider.getBalance(deployer.address);
    console.log("✅ 部署者余额:", ethers.formatEther(balance), "ETH");
}

main().catch((error) => {
    console.error("❌ 错误详情:", error);
    process.exit(1);
});