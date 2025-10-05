/*
在测试网上发行一个图文并茂的 NFT
任务目标
1.使用 Solidity 编写一个符合 ERC721标准的 NFT合约。
2.将图文数据上传到 IPFS，生成元数据链接。
3.将合约部署到以太坊测试网(如Goerli或Sepolia)。
4.铸造 NFT 并在测试网环境中查看。
任务步骤
1.编写 NFT合约
。使用OpenZeppelin的ERC721库编写一个NFT合约:
。合约应包含以下功能:
。构造函数:设置 NFT的名称和符号。
。mintNFT函数:允许用户铸造NFT，并关联元数据链接(tokenURI)
。在 RemixIDE 中编译合约。
2.准备图文数据
。准备一张图片，并将其上传到IPFS(可以使用 Pinata 或其他工具)
。 创建一个JSON文件，描述NFT的属性(如名称、描述、图片链接等)
。将JSON 文件上传到IPFS，获取元数据链接。
。JSON文件参考https://docs.opensea.io/docs/metadata-standards
3.部署合约到测试网
。在Remix IDE 中连接MetaMask，并确保MetaMask 连接到 Goerli或Sepolia 测试网
。部署 NFT合约到测试网，并记录合约地址。
4.铸造 NFT
。使用 mintNFT函数铸造 NFT:
。在recipient 字段中输入你的钱包地址。
。在 tokenURl字段中输入元数据的 IPFS 链接。
。在 MetaMask中确认交易。
5.查看 NFT
。打开OpenSea测试网或Etherscan 测试网.
。连接你的钱包，查看你铸造的 NFT。
*/

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

contract MyNFT is ERC721, ERC721URIStorage {

    using Strings for uint256;

    uint256 private _tokenIdCounter;

    constructor() ERC721("MyFirstNFT", "MFN") {
        _tokenIdCounter++;
    }

    //recipient: 0x46B848559414E3b80142FD0D34A8818b7D540bef
    //tokenURL:  https://gateway.pinata.cloud/ipfs/bafybeiby4vq6y3nse5dteecpi2gltd6n6bpxarn7t4k4ehrldhwqt6oooi
    function mintNFT(address recipient, string memory tokenURL) public {
        uint256 tokenId = _tokenIdCounter;
        _tokenIdCounter++;
        _safeMint(recipient, tokenId, "This is Kim's NFT");
        _setTokenURI(tokenId, tokenURL);
    }

    function tokenURI(uint256 tokenId) public view override(ERC721, ERC721URIStorage) returns (string memory) {
        _requireOwned(tokenId);

        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId) public view virtual override(ERC721, ERC721URIStorage) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}