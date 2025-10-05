const { ethers, deployments, upgrades } = require("hardhat");
const { expect } = require("chai");

// describe("Test deploy", async function () {
//     it("Should be able to deploy", async function () {
//         const Contract = await ethers.getContractFactory("NFTAuction")
//         const contract = await Contract.deploy()
//         await contract.waitForDeployment()

//         await contract.createAuction(
//             100 * 1000,
//             ethers.parseEther("0.000000000000000001"),
//             ethers.ZeroAddress,
//             1
//         )

//         const auction = await contract.auctions(0)

//         console.log(auction)
//     })
// })

// describe("Test upgrade", async function () {
//     it("Should be able to upgrade", async function () {
//         // 1.部署业务合约
//         await deployments.fixture(["deployNFTAuction"])

//         const nftAuctionProxy = await deployments.get("NFTAuctionProxy")

//         const nftAuction = await ethers.getContractAt("NFTAuction", nftAuctionProxy.address)

//         // 2.调用 createAuction 方法创建拍卖
//         await nftAuction.createAuction(
//             100 * 1000,
//             ethers.parseEther("0.01"),
//             ethers.ZeroAddress,
//             1
//         )

//         const auction = await nftAuction.auctions(0)
//         console.log("升级前的拍卖：", auction)

//         const implAddress1 = await upgrades.erc1967.getImplementationAddress(nftAuctionProxy.address)

//         // 3.升级合约
//         await deployments.fixture(["upgradeNFTAuction"])

//         const implAddress2 = await upgrades.erc1967.getImplementationAddress(nftAuctionProxy.address)

//         // 4.读取合约的 auction[0]
//         const auctionV2 = await nftAuction.auctions(0)
//         console.log("升级后的拍卖：", auctionV2)

//         // 调用 testHello 方法
//         const nftAuctionV2 = await ethers.getContractAt("NFTAuctionV2", nftAuctionProxy.address)
//         const hello = await nftAuctionV2.testHello()
//         console.log("hello:", hello)

//         // 5.断言
//         expect(auctionV2.starTime).to.be.equal(auction.starTIme)
//         expect(implAddress1).to.be.not.equal(implAddress2)
//     })
// })

// describe("Test NFT Auction", async function () {
//     it("Auction should be OK", async function () {
//         await deployments.fixture(["deployNFTAuction"]);
//         const nftAuctionProxy = await deployments.get("NFTAuctionProxy");

//         const [signer, buyer] = await ethers.getSigners();

//         // 1. 部署 ERC721 合约
//         const TestERC721 = await ethers.getContractFactory("TestERC721");
//         const testERC721 = await TestERC721.deploy();
//         await testERC721.waitForDeployment();
//         const testERC721Address = await testERC721.getAddress();
//         console.log("testERC721Address: ", testERC721Address);

//         // mint 10个 NFT
//         for (let i = 0; i < 10; i++) {
//             await testERC721.mint(signer.address, i + 1);
//         }

//         const tokenId = 1;

//         // 2. 调用 createAuction 创建拍卖
//         const nftAuction = await ethers.getContractAt("NFTAuction", nftAuctionProxy.address);

//         // 给拍卖合约授权 NFT 转移
//         await testERC721.connect(signer).setApprovalForAll(nftAuctionProxy.address, true);

//         await nftAuction.createAuction(
//             10, // duration
//             ethers.parseEther("0.01"), // startPrice
//             testERC721Address,
//             tokenId
//         );

//         const auction = await nftAuction.auctions(0);
//         console.log("创建拍卖成功: ", auction);

//         // 3. 买家出价
//         await nftAuction.connect(buyer).placeBid(0, { value: ethers.parseEther("0.02") });

//         // 4. 快进时间 10 秒
//         await ethers.provider.send("evm_increaseTime", [10]);
//         await ethers.provider.send("evm_mine", []);

//         // 5. 结束拍卖
//         await nftAuction.endAuction(0);

//         // 验证结果
//         const auctionResult = await nftAuction.auctions(0);
//         console.log("结束拍卖后读取拍卖成功: ", auctionResult);

//         expect(auctionResult.highestBidder).to.equal(buyer.address);
//         expect(auctionResult.highestBid).to.equal(ethers.parseEther("0.02"));

//         // 验证 NFT 所有权
//         const owner = await testERC721.ownerOf(tokenId);
//         console.log("owner: ", owner);
//         expect(owner).to.equal(buyer.address);
//     })
// })

describe("Test contract data feed", async function () {
    this.networkName = "sepolia";
    let nftAuction, ownerSigner;

    // Sepolia 测试网配置（ETH/USD 喂价地址 + 小数位数）
    const SEPOLIA_CONFIG = {
        ETH: {
            tokenAddress: ethers.ZeroAddress, // ETH 的代币地址固定为 0x0
            priceFeedAddress: "0x694AA1769357215DE4FAC081bf1f309aDC325306", // Sepolia ETH/USD 喂价地址
            decimals: 8 // Chainlink 喂价默认 8 位小数（如 1800 USD → 180000000000）
        }
    };

    beforeEach(async function () {
        // 1. 部署合约
        await deployments.fixture(["deployNFTAuction"]);
        const nftAuctionProxy = await deployments.get("NFTAuctionProxy");
        nftAuction = await ethers.getContractAt("NFTAuction", nftAuctionProxy.address);

        // 2. 获取合约 owner 签名者
        const ownerAddress = await nftAuction.owner();
        ownerSigner = await ethers.getSigner(ownerAddress);

        console.log("✅ 合约部署完成");
        console.log("  - Owner 地址：", await ownerSigner.getAddress());
        console.log("  - NFTAuction 代理地址：", await nftAuction.getAddress());
    });

    it("Should set ETH price feed and get latest ETH/USD price", async function () {
        const { tokenAddress, priceFeedAddress, decimals: expectedDecimals } = SEPOLIA_CONFIG.ETH;

        // 3. 调用 setPriceFeed
        try {
            const setTx = await nftAuction.connect(ownerSigner).setPriceFeed(
                tokenAddress,
                priceFeedAddress,
                expectedDecimals,
                { gasLimit: 300000 } // 手动设置 gas，避免网络波动
            );
            await setTx.wait(); // 等待交易上链，确保状态更新
            console.log("\n✅ setPriceFeed 执行成功");
            console.log("  - 交易哈希：", setTx.hash);
            console.log("  - 配置的 ETH 喂价地址：", priceFeedAddress);
        } catch (error) {
            console.error("\n❌ setPriceFeed 调用失败：", error.message.slice(0, 200)); // 截取关键错误信息
            throw error; // 终止测试，避免后续无效执行
        }

        // 4. 调用 getLatestPrice
        try {
            const [rawPrice, actualDecimals] = await nftAuction.getLatestPrice(tokenAddress);
            // 格式化价格：用合约返回的 decimals（而非硬编码 8），更灵活
            const ethUsdPrice = ethers.formatUnits(rawPrice, actualDecimals);

            // 5. 断言验证（确保数据有效）
            expect(actualDecimals).to.equal(expectedDecimals, "喂价小数位数不匹配"); // 校验小数位数
            expect(rawPrice).to.be.gt(0n, "ETH/USD 价格不能为 0"); // 非零校验
            expect(parseFloat(ethUsdPrice)).to.be.gt(1000, "ETH 价格应大于 1000 USD（常识校验）");

            // 打印结果（格式化后更易读）
            console.log("\n✅ 测试成功！ETH/USD 价格信息：");
            console.log(`  - 原始喂价数据：${rawPrice.toString()}`);
            console.log(`  - 喂价小数位数：${actualDecimals}`);
            console.log(`  - 格式化价格：${ethUsdPrice} USD`);
        } catch (error) {
            console.error("\n❌ getLatestPrice 调用失败：", error.message.slice(0, 200));
            throw error;
        }
    });

    // 可选：新增测试 - 未设置喂价时调用 getLatestPrice 应 revert
    // it("should revert if price feed is not set", async function () {
    //     const randomTokenAddress = ethers.Wallet.createRandom().address; // 随机未配置的代币地址
    //     await expect(
    //         nftAuction.getLatestPrice(randomTokenAddress)
    //     ).to.be.revertedWith("Price feed not set"); // 断言 revert 原因
    //     console.log("\n✅ 未设置喂价时调用 getLatestPrice，正确触发 revert");
    // });
});
