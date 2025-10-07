const { expect } = require("chai");
const { ethers } = require("hardhat");

// ------------------------------
// 测试配置（根据需求调整）
// ------------------------------
const TEST_CONFIG = {
    ETH_FEED_DECIMALS: 8,        // 模拟ETH/USD预言机小数位数
    ETH_PRICE: ethers.parseUnits("2000", 8), // 模拟ETH价格：2000 USD
    AUCTION_DURATION: 60,        // 拍卖时长：60秒（本地测试足够）
    ETH_START_PRICE: ethers.parseEther("0.001"), // 起始价：0.001 ETH
    TEST_NFT_NAME: "TestNFT",    // 测试NFT名称
    TEST_NFT_SYMBOL: "TNFT",     // 测试NFT符号
    TEST_TOKEN_ID: 1n,           // 测试NFT的TokenID
};

// ------------------------------
// 模拟预言机合约（本地测试用，替代Chainlink）
// ------------------------------
async function deployMockAggregator(decimals, price) {
    const MockAggregator = await ethers.getContractFactory("MockAggregatorV3");
    const mock = await MockAggregator.deploy(decimals, price);
    await mock.waitForDeployment();
    return mock;
}

// ------------------------------
// 核心测试脚本
// ------------------------------
describe("NFTAuction 全功能测试", function () {
    // 全局超时（本地网络30秒足够，测试网需调整为180000ms）
    this.timeout(30000);

    // 测试角色与合约实例
    let deployer, seller, bidder1, bidder2, randomUser;
    let mockETHFeed, testNFT, factory, auctionInstance;
    let factoryAddr, nftAddr, mockETHFeedAddr;

    // ------------------------------
    // 测试前置准备（一次性部署合约）
    // ------------------------------
    before(async function () {
        // 1. 获取测试账户（Hardhat本地网络默认生成10个账户）
        [deployer, seller, bidder1, bidder2, randomUser] = await ethers.getSigners();
        console.log(`✅ 加载测试账户：\n部署者: ${deployer.address}\n卖家: ${seller.address}\n买家1: ${bidder1.address}\n买家2: ${bidder2.address}`);

        // 2. 部署模拟ETH价格预言机
        mockETHFeed = await deployMockAggregator(
            TEST_CONFIG.ETH_FEED_DECIMALS,
            TEST_CONFIG.ETH_PRICE
        );
        mockETHFeedAddr = await mockETHFeed.getAddress();
        console.log(`✅ 部署模拟ETH预言机：${mockETHFeedAddr}`);

        // 3. 部署测试NFT合约（TestERC721）
        const TestERC721 = await ethers.getContractFactory("TestERC721");
        testNFT = await TestERC721.deploy();
        await testNFT.waitForDeployment();
        nftAddr = await testNFT.getAddress();
        console.log(`✅ 部署测试NFT合约：${nftAddr}`);

        // 4. 铸造测试NFT给卖家（仅owner可铸造，需用deployer授权或直接铸造）
        await testNFT.mint(seller.address, TEST_CONFIG.TEST_TOKEN_ID);
        expect(await testNFT.ownerOf(TEST_CONFIG.TEST_TOKEN_ID)).to.equal(seller.address);
        console.log(`✅ 铸造NFT给卖家，TokenID: ${TEST_CONFIG.TEST_TOKEN_ID}`);

        // 5. 部署工厂合约（NFTAuctionFactory）并初始化
        const NFTAuctionFactory = await ethers.getContractFactory("NFTAuctionFactory");
        factory = await NFTAuctionFactory.deploy();
        await factory.initialize(); // 工厂合约初始化
        factoryAddr = await factory.getAddress();
        console.log(`✅ 部署工厂合约：${factoryAddr}`);

        // 6. 配置ETH的预言机（仅owner可配置）
        await factory.setBidTokenFeed(
            ethers.ZeroAddress, // ETH对应地址0
            mockETHFeedAddr,
            TEST_CONFIG.ETH_FEED_DECIMALS
        );
        const [configuredFeed, configuredDec] = await factory.getBidTokenFeed(ethers.ZeroAddress);
        expect(configuredFeed).to.equal(mockETHFeedAddr);
        expect(configuredDec).to.equal(TEST_CONFIG.ETH_FEED_DECIMALS);
        console.log(`✅ 配置ETH预言机完成`);
    });

    // ------------------------------
    // 1. 工厂合约基础功能测试
    // ------------------------------
    describe("1. NFTAuctionFactory 基础功能", function () {
        it("1.1 初始化后 owner 应为部署者", async function () {
            expect(await factory.owner()).to.equal(deployer.address);
        });

        it("1.2 仅 owner 可配置预言机（普通用户配置失败）", async function () {
            // 普通用户（卖家）尝试配置预言机，应revert
            await expect(
                factory.connect(seller).setBidTokenFeed(
                    ethers.ZeroAddress,
                    mockETHFeedAddr,
                    TEST_CONFIG.ETH_FEED_DECIMALS
                )
            ).to.be.revertedWithCustomError(factory, "OwnableUnauthorizedAccount");
        });

        it("1.3 可正确获取已配置的预言机信息", async function () {
            const [feedAddr, feedDec] = await factory.getBidTokenFeed(ethers.ZeroAddress);
            expect(feedAddr).to.equal(mockETHFeedAddr);
            expect(feedDec).to.equal(TEST_CONFIG.ETH_FEED_DECIMALS);
        });
    });

    // ------------------------------
    // 2. 创建拍卖功能测试
    // ------------------------------
    describe("2. 创建拍卖（ETH竞价）", function () {
        before(async function () {
            // 卖家授权NFT给拍卖合约（创建拍卖前需授权）
            await testNFT.connect(seller).setApprovalForAll(factoryAddr, true);
            console.log(`✅ 卖家授权NFT给工厂合约`);
        });

        it("2.1 NFT所有者（卖家）可创建拍卖", async function () {
            // 监听AuctionCreated事件，验证参数
            const auctionCreatedTx = await factory.connect(seller).createAuction(
                TEST_CONFIG.AUCTION_DURATION,
                TEST_CONFIG.ETH_START_PRICE,
                nftAddr,
                TEST_CONFIG.TEST_TOKEN_ID,
                ethers.ZeroAddress // 用ETH竞价
            );

            // 解析事件
            const receipt = await auctionCreatedTx.wait();
            const auctionCreatedEvent = receipt.events.find(e => e.event === "AuctionCreated");
            const auctionAddr = auctionCreatedEvent.args.auctionAddr;

            // 验证拍卖实例有效
            auctionInstance = await ethers.getContractAt("NFTAuction", auctionAddr);
            expect(await auctionInstance.owner()).to.equal(seller.address); // 拍卖合约owner为卖家

            // 验证工厂映射状态
            const [auctionFromMap, isAuctioning] = await factory.getNFTAuctionInfo(nftAddr, TEST_CONFIG.TEST_TOKEN_ID);
            expect(auctionFromMap).to.equal(auctionAddr);
            expect(isAuctioning).to.be.true;

            // 验证NFT已转移到拍卖合约
            expect(await testNFT.ownerOf(TEST_CONFIG.TEST_TOKEN_ID)).to.equal(auctionAddr);
            console.log(`✅ 创建拍卖实例：${auctionAddr}`);
        });

        it("2.2 非NFT所有者无法创建拍卖", async function () {
            // 买家1（非NFT所有者）尝试创建拍卖，应revert
            await expect(
                factory.connect(bidder1).createAuction(
                    TEST_CONFIG.AUCTION_DURATION,
                    TEST_CONFIG.ETH_START_PRICE,
                    nftAddr,
                    TEST_CONFIG.TEST_TOKEN_ID,
                    ethers.ZeroAddress
                )
            ).to.be.revertedWith("S!=O"); // S!=O: Seller != Owner
        });

        it("2.3 NFT未授权给工厂时无法创建拍卖", async function () {
            // 铸造新NFT给买家1，未授权
            const newTokenId = 2n;
            await testNFT.mint(bidder1.address, newTokenId);
            await expect(
                factory.connect(bidder1).createAuction(
                    TEST_CONFIG.AUCTION_DURATION,
                    TEST_CONFIG.ETH_START_PRICE,
                    nftAddr,
                    newTokenId,
                    ethers.ZeroAddress
                )
            ).to.be.revertedWith("NFT:App"); // NFT:App: 未授权
        });
    });

    // ------------------------------
    // 3. 竞价功能测试（ETH竞价）
    // ------------------------------
    describe("3. 拍卖竞价（ETH）", function () {
        let bidder1InitialBal, bidder2InitialBal;

        before(async function () {
            // 记录买家初始余额（用于验证退款）
            bidder1InitialBal = await ethers.provider.getBalance(bidder1.address);
            bidder2InitialBal = await ethers.provider.getBalance(bidder2.address);
        });

        it("3.1 买家1可提交有效竞价（高于起始价）", async function () {
            const bidAmount = ethers.parseEther("0.002"); // 高于起始价0.001 ETH
            const placeBidTx = await auctionInstance.connect(bidder1).placeBid(bidAmount, {
                value: bidAmount, // ETH竞价需匹配msg.value
            });

            // 验证事件
            const receipt = await placeBidTx.wait();
            const bidPlacedEvent = receipt.events.find(e => e.event === "BidPlaced");
            expect(bidPlacedEvent.args.bidder).to.equal(bidder1.address);
            expect(bidPlacedEvent.args.bidAmount).to.equal(bidAmount);

            // 验证拍卖状态
            const auctionData = await auctionInstance.auction();
            expect(auctionData.highestBidder).to.equal(bidder1.address);
            expect(auctionData.highestBid).to.equal(bidAmount);
        });

        it("3.2 买家2出价更高时，买家1可收到退款", async function () {
            const higherBid = ethers.parseEther("0.003"); // 高于买家1的0.002 ETH
            const gasPrice = await ethers.provider.getGasPrice();

            // 提交更高出价
            const placeBidTx = await auctionInstance.connect(bidder2).placeBid(higherBid, {
                value: higherBid,
                gasPrice: gasPrice,
            });
            const receipt = await placeBidTx.wait();
            const gasUsed = receipt.gasUsed * gasPrice;

            // 验证买家1收到退款（余额 ≈ 初始余额 + 买家1的出价 - 少量Gas损耗）
            const bidder1FinalBal = await ethers.provider.getBalance(bidder1.address);
            const expectedMinBal = bidder1InitialBal + TEST_CONFIG.ETH_START_PRICE - gasUsed;
            expect(bidder1FinalBal).to.be.gt(expectedMinBal);

            // 验证当前最高出价为买家2
            const auctionData = await auctionInstance.auction();
            expect(auctionData.highestBidder).to.equal(bidder2.address);
            expect(auctionData.highestBid).to.equal(higherBid);
        });

        it("3.3 出价低于当前最高时失败", async function () {
            const lowBid = ethers.parseEther("0.0015"); // 低于买家2的0.003 ETH
            await expect(
                auctionInstance.connect(bidder1).placeBid(lowBid, { value: lowBid })
            ).to.be.revertedWith("BidUSD not highest");
        });

        it("3.4 拍卖超时后无法竞价", async function () {
            // 等待拍卖超时（60秒 + 5秒缓冲）
            await new Promise(resolve => setTimeout(resolve, (TEST_CONFIG.AUCTION_DURATION + 5) * 1000));

            const newBid = ethers.parseEther("0.004");
            await expect(
                auctionInstance.connect(bidder1).placeBid(newBid, { value: newBid })
            ).to.be.revertedWith("Auction timeout");
        });
    });

    // ------------------------------
    // 4. 结束拍卖功能测试
    // ------------------------------
    describe("4. 结束拍卖", function () {
        let sellerInitialBal;

        before(async function () {
            // 记录卖家初始余额（用于验证收款）
            sellerInitialBal = await ethers.provider.getBalance(seller.address);
        });

        it("4.1 仅有权限者（卖家/最高出价者/工厂）可结束拍卖", async function () {
            // 随机用户（无权限）结束拍卖，失败
            await expect(
                auctionInstance.connect(randomUser).endAuction()
            ).to.be.revertedWith("No end perm");

            // 卖家（有权限）结束拍卖，成功
            const endAuctionTx = await auctionInstance.connect(seller).endAuction();
            await endAuctionTx.wait();
            console.log(`✅ 卖家成功结束拍卖`);
        });

        it("4.2 结束后NFT转移给最高出价者（买家2）", async function () {
            expect(await testNFT.ownerOf(TEST_CONFIG.TEST_TOKEN_ID)).to.equal(bidder2.address);
        });

        it("4.3 结束后卖家收到最高出价资金", async function () {
            const highestBid = await (await auctionInstance.auction()).highestBid;
            const gasPrice = await ethers.provider.getGasPrice();

            // 验证卖家余额 ≈ 初始余额 + 最高出价 - Gas损耗
            const sellerFinalBal = await ethers.provider.getBalance(seller.address);
            const expectedMinBal = sellerInitialBal + highestBid - (100000n * gasPrice); // 估算Gas
            expect(sellerFinalBal).to.be.gt(expectedMinBal);
        });

        it("4.4 结束后拍卖状态更新，工厂映射清理", async function () {
            // 验证拍卖已结束
            const auctionData = await auctionInstance.auction();
            expect(auctionData.ended).to.be.true;

            // 验证工厂映射清理（isAuctioning变为false）
            const [_, isAuctioning] = await factory.getNFTAuctionInfo(nftAddr, TEST_CONFIG.TEST_TOKEN_ID);
            expect(isAuctioning).to.be.false;
        });

        it("4.5 无出价时结束拍卖，NFT退回卖家", async function () {
            // 铸造新NFT给卖家
            const newTokenId = 3n;
            await testNFT.mint(seller.address, newTokenId);
            await testNFT.connect(seller).setApprovalForAll(factoryAddr, true);

            // 创建无出价的拍卖
            const createTx = await factory.connect(seller).createAuction(
                10, // 短时长10秒
                TEST_CONFIG.ETH_START_PRICE,
                nftAddr,
                newTokenId,
                ethers.ZeroAddress
            );
            const receipt = await createTx.wait();
            const newAuctionAddr = receipt.events.find(e => e.event === "AuctionCreated").args.auctionAddr;
            const newAuction = await ethers.getContractAt("NFTAuction", newAuctionAddr);

            // 等待超时后结束拍卖
            await new Promise(resolve => setTimeout(resolve, 15000));
            await newAuction.connect(seller).endAuction();

            // 验证NFT退回卖家
            expect(await testNFT.ownerOf(newTokenId)).to.equal(seller.address);
        });
    });

    // ------------------------------
    // 5. 异常场景全覆盖
    // ------------------------------
    describe("5. 异常场景验证", function () {
        it("5.1 拍卖已结束后再次结束", async function () {
            await expect(
                auctionInstance.connect(seller).endAuction()
            ).to.be.revertedWith("Auction ended");
        });

        it("5.2 创建拍卖时起始价为0", async function () {
            const newTokenId = 4n;
            await testNFT.mint(seller.address, newTokenId);
            await testNFT.connect(seller).setApprovalForAll(factoryAddr, true);

            await expect(
                factory.connect(seller).createAuction(
                    TEST_CONFIG.AUCTION_DURATION,
                    0n, // 起始价为0
                    nftAddr,
                    newTokenId,
                    ethers.ZeroAddress
                )
            ).to.be.revertedWith("S>0");
        });

        it("5.3 拍卖时长≤10秒", async function () {
            const newTokenId = 5n;
            await testNFT.mint(seller.address, newTokenId);
            await testNFT.connect(seller).setApprovalForAll(factoryAddr, true);

            await expect(
                factory.connect(seller).createAuction(
                    5, // 时长5秒（≤10）
                    TEST_CONFIG.ETH_START_PRICE,
                    nftAddr,
                    newTokenId,
                    ethers.ZeroAddress
                )
            ).to.be.revertedWith("D>10");
        });
    });
});