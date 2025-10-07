# 项目描述
这是一个NFT拍卖合约项目，支持ERC20 和以太坊到美元的价格。


测试网合约部署到测试网
npx hardhat run deploy/01_deploy_nft_auction.js --network sepolia

测试网合约升级
npx hardhat run deploy/02_upgrade_nft_auction.js --network sepolia

本地网测试合约部署
npx hardhat deploy

测试网执行合约测试用例
npx hardhat test --network sepolia
