const { deployments, upgrades, ethers } = require("hardhat");

const fs = require("fs");
const path = require("path");

module.exports = async ({ getNamedAccounts, deployments }) => {
    const { save } = deployments;
    const { deployer } = await getNamedAccounts();

    console.log("部署用户地址：", deployer);
    const nftAuction = await ethers.getContractFactory("NFTAuction");

    // 通过代理合约部署
    const nftAuctionProxy = await upgrades.deployProxy(nftAuction, [], {
        initializer: "initialize",
        kind: "uups",
    });

    // 等待合约部署完成
    await nftAuctionProxy.waitForDeployment();
    const proxyAddress = await nftAuctionProxy.getAddress();
    console.log("代理合约地址：", proxyAddress);
    const implAddress = await upgrades.erc1967.getImplementationAddress(proxyAddress);
    console.log("实现合约地址：", implAddress);

    const storePath = path.resolve(__dirname, "./.cache/proxyNFTAuction.json");

    fs.writeFileSync(
        storePath,
        JSON.stringify({
            proxyAddress,
            implAddress,
            abi: nftAuction.interface.format("json"),
        })
    );

    await save("NFTAuctionProxy", {
        abi: nftAuction.interface.format("json"),
        address: proxyAddress,
        // args: [],
        // log: true,
    })
}

module.exports.tags = ["deployNFTAuction"];