require("@nomicfoundation/hardhat-toolbox");
require('hardhat-deploy');
require("@openzeppelin/hardhat-upgrades");
require("dotenv").config();

module.exports = {
  // 正确：Solidity 配置（version + settings 同级）
  solidity: {
    version: "0.8.20",
    settings: { // 编译器优化配置必须在这里！
      optimizer: {
        enabled: true,
        runs: 1, // 极致体积优化
        details: {
          yul: true, // 启用 Yul 优化（减少 10%-15% 体积）
        }
      }
    }
  },
  networks: {
    sepolia: {
      url: `https://sepolia.infura.io/v3/${process.env.INFURA_API_KEY}`,
      accounts: [
        process.env.PRIVATE_KEY_DEPLOYER,
        process.env.PRIVATE_KEY_SELLER,
        process.env.PRIVATE_KEY_BIDDER1,
        process.env.PRIVATE_KEY_BIDDER2,
        process.env.PRIVATE_KEY_RANDOMUSER
      ].filter(Boolean),
      gas: 3000000,
      timeout: 120000,
    },
    hardhat: {
      allowUnlimitedContractSize: true
    }
  },
  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY // 建议开启，方便验证合约
  }
};