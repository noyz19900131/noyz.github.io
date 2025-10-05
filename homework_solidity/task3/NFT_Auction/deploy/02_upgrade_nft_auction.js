const { ethers, upgrades } = require("hardhat");

const fs = require("fs");
const path = require("path");

module.exports = async ({ getNamedAccounts, deployments }) => {
    const { save } = await deployments;
    const { deployer } = await getNamedAccounts();
    console.log("部署用户地址：", deployer);

    // 读取文件 ./.cache/proxyNftAuction.json
    const storePath = path.resolve(__dirname, "./.cache/proxyNFTAuction.json");
    const storeData = fs.readFileSync(storePath, "utf-8");
    const { proxyAddress } = JSON.parse(storeData);

    // 升级版的业务合约
    const nftAuctionV2 = await ethers.getContractFactory("NFTAuctionV2");

    // 升级代理合约
    const nftAuctionProxyV2 = await upgrades.upgradeProxy(proxyAddress, nftAuctionV2);

    // 等待合约部署完成
    await nftAuctionProxyV2.waitForDeployment();

    const proxyAddressV2 = await nftAuctionProxyV2.getAddress();
    console.log("代理合约地址：", proxyAddressV2);

    const implAddressV2 = await upgrades.erc1967.getImplementationAddress(proxyAddressV2);
    console.log("实现合约地址：", implAddressV2);

    await save("NFTAuctionProxyV2", {
        abi: nftAuctionV2.interface.format("json"),
        address: proxyAddressV2,
        // args: [],
        // log: true,
    })
}

module.exports.tags = ["upgradeNFTAuction"];