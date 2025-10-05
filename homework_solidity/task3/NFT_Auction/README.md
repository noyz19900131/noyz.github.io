# 项目描述
这是一个NFT拍卖合约项目，支持ERC20 和以太坊到美元的价格。


# 测试网合约部署到测试网
npx hardhat run deploy/01_deploy_nft_auction.js --network sepolia

# 测试网合约升级
npx hardhat run deploy/02_upgrade_nft_auction.js --network sepolia

# 本地网测试合约部署
npx hardhat deploy

# 测试网执行合约测试用例
npx hardhat test --network sepolia

# 参考：
部署用户地址： 0x46B848559414E3b80142FD0D34A8818b7D540bef
代理合约地址： 0x4d1a4faDA77DfbF2791D0653D2a46567794f2f7F
实现合约地址： 0xFa5374565f1C36FBb9f10b68CdE4ed1836D7Fe96
